package main

import (
	"fmt"
	"os"

	"github.com/ZiplEix/mail-toolchain/shared/config"
	"github.com/ZiplEix/mail-toolchain/shared/database"
	"github.com/ZiplEix/mail-toolchain/smtp-server/server"
)

func init() {
	if err := config.LoadEnv("../.env"); err != nil {
		panic(fmt.Sprintf("Failed to load environment variables: %v", err))
	}

	if err := database.Init(os.Getenv("POSTGRES_URL")); err != nil {
		panic(fmt.Sprintf("Failed to initialize database connection: %v", err))
	}

	certPath := os.Getenv("CERT_PATH")
	keyPath := os.Getenv("KEY_PATH")

	if err := server.LoadTLSConfig(certPath, keyPath); err != nil {
		panic(fmt.Sprintf("Failed to load TLS configuration: %v", err))
	}
}

func main() {
	server.LunchSMTPServer()
}
