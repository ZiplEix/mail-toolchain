package main

import (
	"fmt"
	"os"

	"github.com/ZiplEix/mail-toolchain/shared/database"
	"github.com/joho/godotenv"
)

// create_user is a command-line tool to create a new user in the database.
// It takes an email and a password as arguments, hashes the password,
// and stores the user in the database.
//
// It must to run at the root of the project with the .env file present.
func main() {
	godotenv.Load(".env")

	dsn := os.Getenv("POSTGRES_URL")
	if dsn == "" {
		fmt.Println("POSTGRES_URL environment variable is not set")
		return
	}

	database.Init(dsn)

	database.MigrateUsersTable()

	if len(os.Args) != 3 {
		fmt.Println("Usage: create_user <email> <password>")
		return
	}

	email := os.Args[1]
	password := os.Args[2]

	if err := database.CreateUser(email, password); err != nil {
		fmt.Println("Failed to create user: " + err.Error())
		return
	}
	println("User created successfully:", email)
}
