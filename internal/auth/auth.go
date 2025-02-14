package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 5)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CheckPasswordHash(password, hash string) error {
	// it's compare HASH and PASSWORD respectively :^3
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}
