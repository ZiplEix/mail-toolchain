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

	if err := database.MigrateMailsTable(); err != nil {
		panic(fmt.Sprintf("Failed to migrate mails table: %v", err))
	}

	if err := database.MigrateUsersTable(); err != nil {
		panic(fmt.Sprintf("Failed to migrate users table: %v", err))
	}

	certPath := os.Getenv("CERT_PATH")
	keyPath := os.Getenv("KEY_PATH")

	if err := server.LoadTLSConfig(certPath, keyPath); err != nil {
		panic(fmt.Sprintf("Failed to load TLS configuration: %v", err))
	}

	if err := server.LoadPrivateKey(os.Getenv("DKIM_PRIVATE_PATH")); err != nil {
		panic(fmt.Sprintf("Failed to load DKIM private key: %v", err))
	}
}

func main() {
	server.LunchSMTPServer()
}
