package crypt

import "golang.org/x/crypto/bcrypt"

func Bcrypt(password string) string {
	result, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(result)
}
func BcryptCheck(encrypted string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(password))
	if err != nil {
		return false
	}
	return true
}
