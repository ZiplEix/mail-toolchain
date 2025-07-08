package server

import (
	"bufio"
	"net"

	"github.com/ZiplEix/mail-toolchain/shared/logger"
)

func StartIMAP(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	session := NewSession(conn)
	session.Send("* OK IMAP4rev1 Service Ready")

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		logger.Event(conn, "C: "+line)
		if err := HandleCommand(session, line); err != nil {
			session.Send("* BAD " + err.Error())
		}
		if session.LoggedOut {
			session.Send("* BYE IMAP4rev1 Service Closing Transmission Channel")
			return
		}
	}
}
