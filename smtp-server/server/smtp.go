package server

import (
	"fmt"
	"net"

	"github.com/ZiplEix/mail-toolchain/shared/logger"
)

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
