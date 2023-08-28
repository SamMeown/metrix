package signer

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type Signer struct {
	key []byte
}

func New(key string) *Signer {
	if key == "" {
		return nil
	}

	return &Signer{
		key: []byte(key),
	}
}

func (s *Signer) GetSignature(content []byte) string {
	h := hmac.New(sha256.New, s.key)
	h.Write(content)
	signature := h.Sum(nil)

	return hex.EncodeToString(signature)
}
