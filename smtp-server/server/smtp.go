package server

import (
	"crypto/tls"
	"fmt"
	"net"

	"github.com/ZiplEix/mail-toolchain/shared/logger"
)

var tlsConfig *tls.Config
var RequireTLSBeforeAuth = true

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
