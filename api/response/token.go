package response

import (
	"altair/pkg/logger"
	"crypto/hmac"
	"crypto/sha512"
)

// TokenInfo - структура ответа получения токена
type TokenInfo struct {
	Domain      string
	Exp         int64
	UserID      uint64
	UserRole    string
	JWT         string
	sigBytes    []byte
	structBytes []byte
}

// Verify - верификация данных
func (t *TokenInfo) Verify(secret string) bool {
	h := hmac.New(sha512.New, []byte(secret))
	_, err := h.Write(t.structBytes)
	if err != nil {
		logger.Warning.Println(err.Error())
		return false
	}

	mac2 := h.Sum(nil)

	return hmac.Equal(t.sigBytes, mac2)
}

// SetSig - установка подписи
func (t *TokenInfo) SetSig(in []byte) {
	t.sigBytes = in
}

// SetStruct - установка структуры
func (t *TokenInfo) SetStruct(in []byte) {
	t.structBytes = in
}
