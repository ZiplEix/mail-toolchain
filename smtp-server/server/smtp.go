package server

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"

	"github.com/ZiplEix/mail-toolchain/shared/database"
	"github.com/ZiplEix/mail-toolchain/shared/logger"
)

var tlsConfig *tls.Config

var errSessionQuit = fmt.Errorf("session quit requested")

func LoadTLSConfig(certPath, keyPath string) error {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return fmt.Errorf("failed to load TLS certificate and key: %v", err)
	}
	tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	return nil
}

func LunchSMTPServer() {
	listener, err := net.Listen("tcp", ":2525")
	if err != nil {
		panic(fmt.Sprintf("Failed to start SMTP server: %v", err))
	}
	logger.Info("SMTP server listening on port 2525")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		session := NewSession(conn)
		go HandleConnection(session)
	}
}

type smtpHandler func(*Session, string) error

var commandHandlers = map[string]smtpHandler{
	"EHLO":       ehlo,
	"HELO":       helo,
	"NOOP":       noop,
	"MAIL FROM:": mailFrom,
	"RCPT TO:":   rcptTo,
	"DATA":       data,
	"RSET":       rset,
	"VRFY":       vrfy,
	"STARTTLS":   startTLS,
	"QUIT":       quit,
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
			for prefix, handler := range commandHandlers {
				if strings.HasPrefix(lineUpper, prefix) {
					err := handler(session, line)
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

func extractEmailAddress(line string, prefix string) string {
	rest := strings.TrimSpace(strings.TrimPrefix(line, prefix))

	start := strings.Index(rest, "<")
	end := strings.Index(rest, ">")

	if start == -1 || end == -1 || end <= start {
		return ""
	}

	return rest[start : end+1]
}
