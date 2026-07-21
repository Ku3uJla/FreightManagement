package features

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(HashedPassword), nil
}

func ComparePassword(password, hashedpassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(hashedpassword))
	if err != nil {
		return false
	}
	return true
}
