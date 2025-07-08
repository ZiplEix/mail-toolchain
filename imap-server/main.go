package main

import (
	"os"

	"github.com/ZiplEix/mail-toolchain/imap-server/server"
	"github.com/ZiplEix/mail-toolchain/shared/database"
	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load("../.env")

	if err := database.Init(os.Getenv("POSTGRES_URL")); err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
}

func main() {
	server.StartIMAP(":1143")
}
