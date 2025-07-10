package server

import (
	"fmt"
	"strings"

	"github.com/ZiplEix/mail-toolchain/shared/database"
	"github.com/ZiplEix/mail-toolchain/shared/logger"
)

type smtpHandler func(*Session, string) error

type commandEntry struct {
	Handler  smtpHandler
	NeedAuth bool
}

var smtpCommands = map[string]commandEntry{
	"EHLO":       {ehlo, false},
	"HELO":       {helo, false},
	"NOOP":       {noop, false},
	"MAIL FROM:": {mailFrom, true},
	"RCPT TO:":   {rcptTo, true},
	"DATA":       {data, true},
	"RSET":       {rset, true},
	"VRFY":       {vrfy, false},
	"STARTTLS":   {startTLS, false},
	"AUTH":       {auth, false},
	"QUIT":       {quit, false},
}

func HandleConnection(session *Session) {
	defer func() {
		if r := recover(); r != nil {
			logger.Event(session.Conn, fmt.Sprintf("Recovered from panic: %v", r))
			session.Close()
		}
	}()
	defer session.Close()

	session.SendLine("220 localhost SMTP ready")
	logger.Event(session.Conn, "Connection opened")

	var dataLines []string

	for {
		line, err := session.Reader.ReadString('\n')
		if err != nil {
			logger.Event(session.Conn, "Connection closed")
			break
		}
		line = strings.TrimRight(line, "\r\n")
		logger.Event(session.Conn, fmt.Sprintf("C: %s", line))

		switch session.Mode {
		case "command":
			lineUpper := strings.ToUpper(line)

			handled := false
			for prefix, entry := range smtpCommands {
				if strings.HasPrefix(lineUpper, prefix) {
					if entry.NeedAuth && !session.Authenticated {
						session.SendError(530, "Authentication required")
						handled = true
						break
					}

					err := entry.Handler(session, line)
					if err == errSessionQuit {
						return
					}
					handled = true
					break
				}
			}

			if !handled {
				session.SendError(500, "Unrecognized command")
			}

		case "data":
			if line == "." {
				err := database.SaveMail(session.From, session.ToList, dataLines)
				if err != nil {
					session.SendError(550, "Error saving message to database")
					logger.Event(session.Conn, fmt.Sprintf("Error saving mail from %s to %v: %v", session.From, session.ToList, err))
				} else {
					session.SendLine("250 OK: message accepted")
					logger.Event(session.Conn, fmt.Sprintf("Mail saved from %s to %v", session.From, session.ToList))
				}
				session.Mode = "command"
				session.From = ""
				session.ToList = nil
				dataLines = nil
			} else {
				dataLines = append(dataLines, line)
			}
		}
	}
}
