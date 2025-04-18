package token

import (
	"JwtAuth/pkg/generator"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"log"
	"time"
)

func GenerateAccessToken(guid, ip string, expiration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"guid": guid,
		"ip":   ip,
		"exp":  time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateRefreshToken() (string, string, error) {
	refreshToken := generator.GenerateSecureToken(32)
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}
	//
	//encodedToken := base64.StdEncoding.EncodeToString([]byte(refreshToken))
	log.Println(refreshToken, hashedToken)
	return refreshToken, string(hashedToken), nil
}
