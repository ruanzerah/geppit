package utils

import "golang.org/x/crypto/bcrypt"

// It will hash with default the cost "10", It returns an empty string if there's an error
func Hash(input string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// returns true if s == i , otherwise returns false
func CompareHash(stored, input string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(stored), []byte(input))
	return err == nil
}
