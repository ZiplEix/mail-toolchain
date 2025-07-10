package server

import (
	"bufio"
	"net"
	"strings"

	"github.com/ZiplEix/mail-toolchain/shared/logger"
)

type IMAPState int

const (
	StateNotAuthenticated IMAPState = iota
	StateAuthenticated
	StateSelected
)

type Session struct {
	Conn       net.Conn
	Reader     *bufio.Reader
	Writer     *bufio.Writer
	State      IMAPState
	SelectedMB string
	LoggedOut  bool
	Username   string
}

func NewSession(conn net.Conn) *Session {
	return &Session{
		Conn:   conn,
		Reader: bufio.NewReader(conn),
		Writer: bufio.NewWriter(conn),
		State:  StateNotAuthenticated,
	}
}

func (s *Session) Send(msg string) {
	s.Writer.WriteString(msg + "\r\n")
	s.Writer.Flush()

	logger.Debugf("Sent to %s: %s", s.Conn.RemoteAddr().String(), msg)
}

func (s *Session) ReadLine() (string, error) {
	line, err := s.Reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	line = strings.TrimRight(line, "\r\n")
	return line, nil
}
