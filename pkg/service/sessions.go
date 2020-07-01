package service

import (
	"altair/api/response"
	"altair/pkg/manager"
	"altair/server"
	"altair/storage"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"hash"
	"strings"
	"time"
)

// NewSessionService - фабрика, создает объект сессии
func NewSessionService() *SessionService {
	return new(SessionService)
}

// SessionService - структура сессии
type SessionService struct{}

// GetSessionByRefreshToken - получить сессию по рефреш-токену. Тут отдаем копию.
func (ss SessionService) GetSessionByRefreshToken(refreshToken string) (storage.Session, error) {
	session := storage.Session{}
	err := server.Db.Where("refresh_token = ?", refreshToken).First(&session).Error

	return session, err
}

// GetSessionByUserID - получить сессию относительно ID пользователя
func (ss SessionService) GetSessionByUserID(userID uint64) ([]*storage.Session, error) {
	sessions := make([]*storage.Session, 0)
	err := server.Db.Order("expires_in desc").Where("user_id = ?", userID).Find(&sessions).Error

	return sessions, err
}

// Create - создать сессию
func (ss SessionService) Create(session *storage.Session, tx *gorm.DB) error {
	if !server.Db.NewRecord(session) {
		return manager.ErrNotCreateNewSession
	}

	if tx == nil {
		tx = server.Db
	}

	err := tx.Create(session).Error

	return err
}

// Update - изменить сессию
func (ss SessionService) Update(session *storage.Session, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db
	}

	err := tx.Save(session).Error

	return err
}

// Delete - удалить сессию
func (ss SessionService) Delete(sessionID uint64, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db
	}
	if err := tx.Delete(storage.Session{}, "session_id = ?", sessionID).Error; err != nil {
		return err
	}

	return nil
}

// DeleteAllByUserID - удалить все сессию относительно ID пользователя
func (ss SessionService) DeleteAllByUserID(userID uint64, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db
	}
	if err := tx.Delete(storage.Session{}, "user_id = ?", userID).Error; err != nil {
		return err
	}

	return nil
}

// GenerateAccessToken - сгенерировать аксес-токен
func (ss SessionService) GenerateAccessToken(userID uint64, secret string, userRole string) (response.TokenInfo, error) {
	var h hash.Hash
	tokenInfo := response.TokenInfo{
		Domain:   manager.CookieDomain,
		Exp:      time.Now().Add(time.Second * manager.AccessTokenTimeSecond).Unix(),
		UserID:   userID,
		UserRole: userRole,
	}

	structBytes, err := json.Marshal(tokenInfo)
	if err != nil {
		return tokenInfo, err
	}

	h = hmac.New(sha512.New, []byte(secret))
	if _, err := h.Write(structBytes); err != nil {
		return tokenInfo, err
	}

	sigBytes := h.Sum(nil)
	structTransport := base64.RawURLEncoding.EncodeToString(structBytes)
	sigTransport := base64.RawURLEncoding.EncodeToString(sigBytes)

	tokenInfo.JWT = fmt.Sprintf("%s.%s", structTransport, sigTransport)

	return tokenInfo, nil
}

// ParseAccessToken - разобрать на детали аксес-токен
func (ss SessionService) ParseAccessToken(tokenStr string) (response.TokenInfo, error) {
	tokenInfo := response.TokenInfo{}

	parts := strings.Split(tokenStr, ".")
	if len(parts) != 2 {
		return tokenInfo, manager.ErrNotCorrectClaimsAndSig
	}

	structTransport := parts[0]
	sigTransport := parts[1]

	structBytes, err := base64.RawURLEncoding.DecodeString(structTransport)
	if err != nil {
		return tokenInfo, err
	}

	sigBytes, err := base64.RawURLEncoding.DecodeString(sigTransport)
	if err != nil {
		return tokenInfo, err
	}

	err = json.Unmarshal(structBytes, &tokenInfo)
	if err != nil {
		return tokenInfo, err
	}

	tokenInfo.JWT = tokenStr
	tokenInfo.SetStruct(structBytes)
	tokenInfo.SetSig(sigBytes)

	return tokenInfo, nil
}

// ReloadTokens - перевыпустить токены
func (ss SessionService) ReloadTokens(userID uint64, tokenPassword string, userRole string, c *gin.Context) (response.TokenInfo, *storage.Session, int, error) {
	var session = new(storage.Session)
	var err error

	tokenInfo, err := ss.GenerateAccessToken(userID, tokenPassword, userRole)
	if err != nil {
		return tokenInfo, session, 500, err
	}

	session.UserID = userID
	session.RefreshToken = manager.RandASCII(32)
	session.ExpiresIn = time.Now().Add(time.Second * manager.RefreshTokenTimeSecond)
	session.IP = c.ClientIP()
	session.UserAgent = base64.RawURLEncoding.EncodeToString([]byte(c.GetHeader("User-Agent")))
	session.Fingerprint = ""

	// прежде чем создать сессию нужно проверить на их кол-во. Буффер = 5.
	if sessionsOld, err1 := ss.GetSessionByUserID(userID); err1 != nil {
		return tokenInfo, session, 500, err1

	} else if len(sessionsOld) > manager.SessionLimit {
		tx := server.Db.Begin()
		for _, v := range sessionsOld[manager.SessionLimit:] {
			if err2 := ss.Delete(v.SessionID, tx); err2 != nil {
				tx.Rollback()
				return tokenInfo, session, 500, err2
			}
		}
		tx.Commit()
	}

	if err := ss.Create(session, nil); err != nil {
		return tokenInfo, session, 500, err
	}

	return tokenInfo, session, 200, nil
}

// private -------------------------------------------------------------------------------------------------------------
