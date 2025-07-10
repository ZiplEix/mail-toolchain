package server

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"strings"

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
	logger.Event(s.Conn, "TLS connection established")
	return nil
}
