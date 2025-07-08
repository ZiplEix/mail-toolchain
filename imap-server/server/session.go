package server

import (
	"bufio"
	"net"
)

type IMAPState int

const (
	StateNotAuthenticated IMAPState = iota
	StateAuthenticated
	StateSelected
)

type Session struct {
	Conn       net.Conn
	Writer     *bufio.Writer
	State      IMAPState
	SelectedMB string
	LoggedOut  bool
	Username   string
}

func NewSession(conn net.Conn) *Session {
	return &Session{
		Conn:   conn,
		Writer: bufio.NewWriter(conn),
		State:  StateNotAuthenticated,
	}
}

func (s *Session) Send(msg string) {
	s.Writer.WriteString(msg + "\r\n")
	s.Writer.Flush()

	// fmt.Println("Sent:", msg)
}
