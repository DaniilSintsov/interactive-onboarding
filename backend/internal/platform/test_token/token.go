package test_token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

const rawTokenSize = 32

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Generate() (string, []byte, error) {
	randomBytes := make([]byte, rawTokenSize)

	if _, err := rand.Read(randomBytes); err != nil {
		return "", nil, err
	}

	rawToken := base64.RawURLEncoding.EncodeToString(randomBytes)
	hash := s.Hash(rawToken)

	return rawToken, hash, nil
}

func (s *Service) Hash(rawToken string) []byte {
	hash := sha256.Sum256([]byte(rawToken))
	return hash[:]
}
