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

func (s *Signer) EqualSignatures(sign1, sign2 string) bool {
	return hmac.Equal([]byte(sign1), []byte(sign2))
}

func (s *Signer) ValidateSignature(sign string, content []byte) bool {
	contentSignature := s.GetSignature(content)
	return s.EqualSignatures(contentSignature, sign)
}
