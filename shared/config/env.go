package config

import "github.com/joho/godotenv"

func checkEnv() error {
	return nil
}

func LoadEnv(envPath string) error {
	_ = godotenv.Load(envPath)

	err := checkEnv()
	if err != nil {
		return err
	}

	return nil
}
