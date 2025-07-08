package main

import (
	"fmt"
	"os"

	"github.com/ZiplEix/mail-toolchain/shared/database"
	"github.com/joho/godotenv"
)

// check_user is a command-line tool to check if a user exists in the database
// and if the provided password matches the stored password hash.
// It takes an email and a password as arguments and checks the user credentials.
//
// It must be run at the root of the project with the .env file present.
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

	ok, err := database.CheckUserPassword(email, password)
	if err != nil {
		fmt.Println("Failed to check user password: " + err.Error())
		return
	}
	if !ok {
		fmt.Println("Invalid email or password")
		return
	}
	println("User authenticated successfully:", email)
}
