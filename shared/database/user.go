package database

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = Pool.Exec(context.Background(), `INSERT INTO users (email, password_hash) VALUES ($1, $2)`, email, string(hash))
	return err
}

func CheckUserPassword(email, password string) (bool, error) {
	var hash string
	err := Pool.QueryRow(context.Background(), `SELECT password_hash FROM users WHERE email = $1`, email).Scan(&hash)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return false, nil
		}
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil, err
}
