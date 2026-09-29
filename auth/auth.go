package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nxrmqlly/district/store"
)

const (
	SessionLifetime = 30 * 24 * time.Hour
)

const (
	argonIterations  = 3
	argonMemory      = 64 * 1024 // 64 MiB
	argonParallelism = 4
	argonKeyLen      = 32
	argonSaltLen     = 16
)

type Service struct {
	queries *store.Queries
}

func New(queries *store.Queries) *Service {
	return &Service{queries: queries}
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

func HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, &argon2id.Params{
		Iterations:  argonIterations,
		Memory:      argonMemory,
		Parallelism: argonParallelism,
		SaltLength:  argonSaltLen,
		KeyLength:   argonKeyLen,
	})
}

func VerifyPassword(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

type AuthSession struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Username  string
	CreatedAt time.Time
	ExpiresAt time.Time
	LastUsed  *time.Time
}

type AuthUser struct {
	ID       uuid.UUID
	Username string
	Email    string
}

func (s *Service) CreateSession(ctx context.Context, userID uuid.UUID) (string, error) {
	sessionToken, err := randomToken()
	if err != nil {
		return "", err
	}

	tokHash := sha256sum(sessionToken)

	_, err = s.queries.NewSession(ctx, store.NewSessionParams{
		UserID:    userID,
		TokenHash: tokHash,
		ExpiresAt: time.Now().Add(SessionLifetime),
	})

	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(sessionToken), nil
}

func (s *Service) GetSession(ctx context.Context, tokenB64 string) (*AuthSession, error) {
	raw, err := base64.RawURLEncoding.DecodeString(tokenB64)
	if err != nil {
		return nil, err
	}

	hashed := sha256sum([]byte(raw))
	se, err := s.queries.GetSessionByTokenHash(ctx, hashed)
	if err != nil {
		return nil, err
	}

	return &AuthSession{
		ID:        se.ID,
		UserID:    se.UserID,
		Username:  se.SessionUsername,
		CreatedAt: se.CreatedAt,
		ExpiresAt: se.ExpiresAt,
		LastUsed:  se.LastUsed,
	}, nil
}

func (s *Service) RegisterUser(ctx context.Context, email, username, passwd string) (*AuthUser, error) {
	passwdHash, err := HashPassword(passwd)
	if err != nil {
		return nil, err
	}

	user, err := s.queries.NewUser(ctx, store.NewUserParams{
		Username:   username,
		Email:      email,
		PasswdHash: passwdHash,
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "users_username_key", "users_username_lower_idx":
				return nil, ErrUsernameTaken
			case "users_email_key":
				return nil, ErrEmailInUse
			}
		}
		return nil, err
	}

	return &AuthUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (s *Service) Authenticate(ctx context.Context, login, password string) (*AuthUser, error) {
	user, err := s.queries.GetUserByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	ok, err := VerifyPassword(password, user.PasswdHash)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrInvalidCredentials
	}

	// * creds matched

	return &AuthUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (s *Service) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	return s.queries.RevokeSession(ctx, sessionID)
}
