package keygen

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Generate() (string, error) {
	randomBytes := make([]byte, 32)

	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("project key generator: %w", err)
	}

	return "pk_" + base64.RawURLEncoding.EncodeToString(randomBytes), nil
}
