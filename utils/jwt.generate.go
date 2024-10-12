package utils

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

var JWTSecret []byte

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	JWTSecret = []byte(os.Getenv("JWT_SECRET"))
	if len(JWTSecret) == 0 {
		log.Fatal("JWT_SECRET not set in .env")
	}
}
func GenerateJWT(userId string, email string, role string) (string, error) {

	// Set token expiration time
	expirationtime := time.Now().Add(24 * time.Hour)

	//Create jwt class which contains userId and email and expiration time
	claim := jwt.MapClaims{
		"userId": userId,
		"email":  email,
		"role":   role,
		"exp":    expirationtime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", err
	}
	return tokenString, err

}
