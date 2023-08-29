package token

import (
	"os"

	"github.com/joho/godotenv"
)

var secret []byte

func init() {
	load()
}

func load() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}

	secret = []byte(os.Getenv("SECRET"))
}

func JWTSecret() []byte {
	if len(secret) == 0 {
		load()
	}

	return secret
}
