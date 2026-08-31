package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"

	"github.com/google/uuid"
	iamv1 "github.com/lucas-woo/cloud-drive/api/iam/v1"
	"github.com/lucas-woo/cloud-drive/internal/dto"
	"google.golang.org/protobuf/types/known/timestamppb"
)


func GenerateApiSecret() (plainTextSecret string, dbHash []byte, err error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", nil, errors.New("failed to generate secure bytes")
	}

	plainTextSecret = hex.EncodeToString(bytes)

	hash := sha256.Sum256(bytes)
	
	return plainTextSecret, hash[:], nil
}


func CompareSecret(incomingSecret string, dbHash []byte) bool {
	incomingHash := sha256.Sum256([]byte(incomingSecret))
	return subtle.ConstantTimeCompare(incomingHash[:], dbHash) == 1
}

func ConvertApiKeysToResponse(apiKeys []*dto.ApiKey) []*iamv1.ApiKey {
	result := make([]*iamv1.ApiKey, 0, len(apiKeys))

	for _, apiKey := range apiKeys {
		if apiKey == nil {
			continue
		}

		result = append(result, &iamv1.ApiKey{
			ApiKey:    apiKey.ApiKey,
			KeyName:   apiKey.Name,
			CreatedAt: timestamppb.New(apiKey.CreatedAt),
			IsActive: apiKey.IsActive,
		})
	}

	return result
}


func GenerateApiKey() (uuid.UUID, error) {
	return uuid.NewV7() 
}