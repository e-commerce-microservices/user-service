package service

import "golang.org/x/crypto/bcrypt"

// createHashedPassword generate new hashed password by bcrypt
func createHashedPassword(password string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(hashed)
}

// compareHashedPassword return true if password and hashedPassword are equal
func compareHashedPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
