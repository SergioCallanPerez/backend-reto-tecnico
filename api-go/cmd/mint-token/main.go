package main

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

const tokenLifetime = 90 * 24 * time.Hour

func main() {
	_ = godotenv.Load()

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		fmt.Fprintln(os.Stderr, "falta JWT_SECRET en el entorno (o en .env)")
		os.Exit(1)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(tokenLifetime).Unix(),
	})

	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Fprintln(os.Stderr, "no se pudo firmar el token:", err)
		os.Exit(1)
	}

	fmt.Println(signed)
}
