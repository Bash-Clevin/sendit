package main

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
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

func getenv(env string, def string) string {
	val := os.Getenv(env)
	if len(val) == 0 {
		return def
	}
	return val
}

type FileHeader struct {
	Filesize int64
}

func parseFileHeader(r *bufio.Reader) (*FileHeader, error) {
	raw, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}
	headerParts := strings.Split(raw, " ")
	if len(headerParts) != 3 {
		return nil, fmt.Errorf("invalid header")
	}

	filesizeStr := headerParts[1]
	filesize, _ := strconv.Atoi(filesizeStr)

	if filesize > maxFileSize {
		return nil, fmt.Errorf("max filesize exceeded: %d > %d", filesize, maxFileSize)
	}

	return &FileHeader{
		Filesize: int64(filesize),
	}, nil
}
