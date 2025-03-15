package auth

import "golang.org/x/crypto/bcrypt"


func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		println("Failed to hash password")
		return "", err
	}

	return string(hash), nil
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(hash))
	if err != nil {
		println("Incorrect password, did not match hash!")
		return err
	}

	return nil
}