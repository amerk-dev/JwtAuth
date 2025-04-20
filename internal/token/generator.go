package token

import (
	"JwtAuth/pkg/generator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"log"
	"time"
)

func GenerateAccessToken(userGuid string, ip string, expiration time.Duration) (string, string, error) {
	jti := uuid.New().String()
	now := time.Now()

	claims := jwt.MapClaims{
		"sub": userGuid,
		"jti": jti,
		"iat": now.Unix(),
		"exp": now.Add(expiration).Unix(),
		"ip":  ip,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	signedToken, err := token.SignedString([]byte("12345678"))
	if err != nil {
		return "", "", err
	}

	return signedToken, jti, nil
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
