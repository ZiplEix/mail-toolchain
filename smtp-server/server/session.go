package server

import (
	"bufio"
	"fmt"
	"net"

	"github.com/ZiplEix/mail-toolchain/shared/logger"
)

type Session struct {
	Conn   net.Conn
	Reader *bufio.Reader
	Writer *bufio.Writer

	From   string
	ToList []string
	Mode   string // "command" or "data"

	Username      string
	Authenticated bool

	TLS bool
}

func NewSession(conn net.Conn) *Session {
	return &Session{
		Conn:          conn,
		Reader:        bufio.NewReader(conn),
		Writer:        bufio.NewWriter(conn),
		From:          "",
		ToList:        nil,
		Mode:          "command",
		Username:      "",
		Authenticated: false,
		TLS:           false,
	}
}

func (s *Session) SendLine(line string) error {
	_, err := s.Writer.WriteString(line + "\r\n")
	if err != nil {
		return err
	}
	logger.Debug("S: " + line)
	return s.Writer.Flush()
}

func (s *Session) SendError(code int, msg string) error {
	resp := fmt.Sprintf("%d %s", code, msg)
	if err := s.SendLine(resp); err != nil {
		return err
	}
	return nil
}

func (s *Session) ReadLine() (string, error) {
	line, err := s.Reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return line[:len(line)-2], nil // Remove \r\n
}

func (s *Session) Close() error {
	if err := s.Conn.Close(); err != nil {
		return err
	}
	return nil
}
