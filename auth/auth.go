package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/google/uuid"

	"github.com/nxrmqlly/district/store"
)

const (
	SessionLifetime = 30 * 24 * time.Hour
)

type Service struct {
	queries *store.Queries
}

func randomToken() ([]byte, error) {
	token := make([]byte, 32)

	if _, err := rand.Read(token); err != nil {
		return nil, err
	}

	return token, nil
}

func sha256sum(t []byte) []byte {
	sha := sha256.Sum256(t)
	return sha[:]
}

type AuthSession struct {
	UserID    uuid.UUID
	CSRFHash  []byte
	CreatedAt time.Time
	ExpiresAt time.Time
	LastUsed  *time.Time
}

func (s *Service) CreateSession(ctx context.Context, userID uuid.UUID) (string, error) {
	sessionToken, err := randomToken()
	if err != nil {
		return "", err
	}

	csrfToken, err := randomToken()
	if err != nil {
		return "", err
	}

	tokHash := sha256sum(sessionToken)
	csrfHash := sha256sum(csrfToken)

	_, err = s.queries.NewSession(ctx, store.NewSessionParams{
		UserID:    userID,
		TokenHash: tokHash,
		CsrfHash:  csrfHash,
		ExpiresAt: time.Now().Add(SessionLifetime),
	})

	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(sessionToken), nil
}

func (s *Service) GetSession(ctx context.Context, token string) (*AuthSession, error) {
	tokHash := sha256sum([]byte(token))
	se, err := s.queries.GetSessionByTokenHash(ctx, tokHash)
	if err != nil {
		return nil, err
	}

	return &AuthSession{
		UserID:    se.UserID,
		CSRFHash:  se.CsrfHash,
		CreatedAt: se.CreatedAt,
		ExpiresAt: se.ExpiresAt,
		LastUsed:  se.LastUsed,
	}, nil
}
