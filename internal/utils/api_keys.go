package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateApiKey() (string, error) {
	temp := make([]byte, 16)

	if _, err := rand.Read(temp); err != nil {
		return "", nil;
	}
	return base64.URLEncoding.EncodeToString(temp), nil	
}

func GenerateApiSecret() (string, error) {
	temp := make([]byte, 16)

	if _, err := rand.Read(temp); err != nil {
		return "", nil;
	}
	return base64.URLEncoding.EncodeToString(temp), nil
}

