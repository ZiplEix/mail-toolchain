package server

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/ZiplEix/mail-toolchain/shared/database"
	"github.com/ZiplEix/mail-toolchain/shared/logger"
	"github.com/ZiplEix/mail-toolchain/smtp-server/server/internal"
)

func noop(s *Session, line string) error {
	s.SendLine("250 OK")
	return nil
}

func ehlo(s *Session, line string) error {
	s.SendLine("250-localhost Hello")
	s.SendLine("250-PIPELINING")
	s.SendLine("250-SIZE 35882577")
	s.SendLine("250-STARTTLS")

	if s.TLS {
		s.SendLine("250-AUTH PLAIN LOGIN")
	}

	s.SendLine("250 HELP")
	return nil
}

func helo(s *Session, line string) error {
	s.SendLine("250 localhost Hello")
	return nil
}

func mailFrom(s *Session, line string) error {
	addr := extractEmailAddress(line, "MAIL FROM:")
	if !internal.IsValidEmail(addr) {
		s.SendError(550, "Invalid sender address")
		return nil
	}
	s.From = addr
	s.SendLine("250 OK")
	return nil
}

func rcptTo(s *Session, line string) error {
	addr := extractEmailAddress(line, "RCPT TO:")
	if !internal.IsValidEmail(addr) {
		s.SendError(550, "Invalid recipient address")
		return nil
	}
	s.ToList = append(s.ToList, addr)
	s.SendLine("250 OK")
	return nil
}

func data(s *Session, line string) error {
	if s.From == "" || len(s.ToList) == 0 {
		s.SendError(503, "Bad sequence: MAIL FROM and RCPT TO required before DATA")
		return nil
	}
	s.SendLine("354 End data with <CR><LF>.<CR><LF>")
	s.Mode = "data"
	return nil
}

func rset(s *Session, line string) error {
	s.From = ""
	s.ToList = nil
	s.Mode = "command"
	s.SendLine("250 OK: Reset state")
	return nil
}

func vrfy(s *Session, line string) error {
	arg := strings.TrimSpace(strings.TrimPrefix(line, "VRFY"))
	if internal.IsValidEmail(arg) {
		s.SendLine("250 User exists")
	} else {
		s.SendError(550, "No such user")
	}
	return nil
}

func quit(s *Session, line string) error {
	s.SendLine("221 Bye")
	s.Close()
	logger.Event(s.Conn, "Connection closed via QUIT")
	return errSessionQuit
}

func startTLS(s *Session, line string) error {
	if tlsConfig == nil {
		s.SendError(454, "TLS not available")
	}
	s.SendLine("220 Ready to start TLS")
	tlsConn := tls.Server(s.Conn, tlsConfig)
	err := tlsConn.Handshake()
	if err != nil {
		logger.Event(s.Conn, fmt.Sprintf("TLS handshake failed: %v", err))
	}
	s.Conn = tlsConn
	s.Reader = bufio.NewReader(tlsConn)
	s.Writer = bufio.NewWriter(tlsConn)
	s.TLS = true
	logger.Event(s.Conn, "TLS connection established")
	return nil
}

func auth(s *Session, line string) error {
	if s.Authenticated {
		s.SendError(503, "Already authenticated")
		return nil
	}

	if RequireTLSBeforeAuth && !s.TLS {
		s.SendError(538, "Encryption required for authentication")
		return nil
	}

	args := strings.SplitN(line, " ", 3)
	if len(args) < 2 {
		s.SendError(501, "Syntax: AUTH <mechanism> [initial-response]")
		return nil
	}

	mechanism := strings.ToUpper(args[1])

	switch mechanism {
	case "PLAIN":
		var b64data string
		if len(args) == 3 {
			b64data = args[2]
		} else {
			s.SendLine("334") // Awaiting base64 blob
			var err error
			b64data, err = s.ReadLine()
			if err != nil {
				return err
			}
		}

		decoded, err := base64.StdEncoding.DecodeString(b64data)
		if err != nil {
			s.SendError(501, "Invalid base64 encoding")
			return nil
		}

		// Expect: \0username\0password
		parts := strings.Split(string(decoded), "\x00")
		if len(parts) != 3 {
			s.SendError(501, "Invalid PLAIN auth format")
			return nil
		}
		username := parts[1]
		password := parts[2]

		return tryAuthenticate(s, username, password)

	case "LOGIN":
		s.SendLine("334 " + base64.StdEncoding.EncodeToString([]byte("Username:")))
		usernameB64, err := s.ReadLine()
		if err != nil {
			return err
		}
		usernameBytes, err := base64.StdEncoding.DecodeString(usernameB64)
		if err != nil {
			s.SendError(501, "Invalid base64 in username")
			return nil
		}
		username := string(usernameBytes)

		s.SendLine("334 " + base64.StdEncoding.EncodeToString([]byte("Password:")))
		passwordB64, err := s.ReadLine()
		if err != nil {
			return err
		}
		passwordBytes, err := base64.StdEncoding.DecodeString(passwordB64)
		if err != nil {
			s.SendError(501, "Invalid base64 in password")
			return nil
		}
		password := string(passwordBytes)

		return tryAuthenticate(s, username, password)

	default:
		s.SendError(504, "Unsupported authentication mechanism")
		return nil
	}
}

func tryAuthenticate(s *Session, username, password string) error {
	ok, err := database.CheckUserPassword(username, password)
	if err != nil {
		s.SendError(454, "Temporary authentication failure")
		return nil
	}
	if !ok {
		s.SendError(535, "Authentication failed")
		return nil
	}

	s.Username = username
	s.Authenticated = true
	s.SendLine("235 Authentication successful")
	return nil
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
