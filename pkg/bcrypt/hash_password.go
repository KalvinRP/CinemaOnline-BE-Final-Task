package bcrypt

import "golang.org/x/crypto/bcrypt"

func HashPass(Password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(Password), 10)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckHash(input string, hashpass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashpass), []byte(input))
	return err == nil
}
