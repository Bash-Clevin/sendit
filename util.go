package main

import (
	"crypto/rand"
	"encoding/base64"
	"regexp"
)

func validateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(email)
}

func generateRandomString(length int) string {
	randomBytes := make([]byte, length)
	rand.Read(randomBytes)
	return base64.URLEncoding.EncodeToString(randomBytes)
}
