package main

import (
	"os"

	mailclient "github.com/ZiplEix/mail-toolchain/smtp-client"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("../.env")

	mailclient.Setup(mailclient.SMTPConfig{
		Host:     "smtp.gmail.com",
		Port:     587,
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
	})
	// mailclient.Setup(mailclient.SMTPConfig{
	// 	Host: "localhost",
	// 	Port: 2525,
	// })

	err := mailclient.SendMail(
		os.Getenv("USERNAME"),
		[]string{os.Getenv("USERNAME")},
		"Test Subject",
		"This is a test email body.",
	)
	if err != nil {
		panic(err)
	}
	println("Email sent successfully!")
}
