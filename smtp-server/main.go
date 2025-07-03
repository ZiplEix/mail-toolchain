package main

import (
	"fmt"
	"os"

	"github.com/ZiplEix/mail-toolchain/smtp-server/db"
	"github.com/ZiplEix/mail-toolchain/smtp-server/server"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../.env")

	err := db.InitDB()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize database: %v", err))
	}

	certPath := os.Getenv("CERT_PATH")
	keyPath := os.Getenv("KEY_PATH")

	if certPath == "" || keyPath == "" {
		panic("CERT_PATH or KEY_PATH environment variable is not set")
	}

	err = server.LoadTLSConfig(certPath, keyPath)
	if err != nil {
		panic(fmt.Sprintf("Failed to load TLS configuration: %v", err))
	}
}

func main() {
	server.LunchSMTPServer()
}
