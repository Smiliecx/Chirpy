package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)


func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		println("Failed to hash password")
		return "", err
	}

	return string(hash), nil
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		println("Incorrect password, did not match hash!")
		return err
	}

	return nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),  
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)), 
		Subject:   userID.String(),  
	})

	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	tokenStruct, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	if claims, ok := tokenStruct.Claims.(*jwt.RegisteredClaims); ok && tokenStruct.Valid {
		subject, err := claims.GetSubject()
		if err != nil {
			return uuid.Nil, err
		}

		userID, err := uuid.Parse(subject)
		if err != nil {
			return uuid.Nil, err
		}

		return userID, nil
	}

	return uuid.Nil, fmt.Errorf("invalid token")
}

func GetBearerToken(headers http.Header) (string, error) {
	auth_header := headers.Get("Authorization")

	if auth_header == "" {
		return "", fmt.Errorf("invalid header")
	}

	if strings.HasPrefix(auth_header, "Bearer ") {
		auth_header = auth_header[7:]
	} else {
		return "", fmt.Errorf("no bearer token found")
	}

	return auth_header, nil
}

func MakeRefreshToken() (string, error) {
	key := make([]byte, 32)
	rand.Read(key)
	as_string := hex.EncodeToString(key)
	return as_string, nil
}