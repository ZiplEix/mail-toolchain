package main

import (
	"fmt"

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

	err = server.LoadTLSConfig("../cert.pem", "../key.pem")
	if err != nil {
		panic(fmt.Sprintf("Failed to load TLS configuration: %v", err))
	}
}

func main() {
	server.LunchSMTPServer()
}
