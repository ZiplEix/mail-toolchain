package server

import (
	"crypto/tls"
	"fmt"
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
