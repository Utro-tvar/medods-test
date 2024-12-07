package service

import (
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/Utro-tvar/medods-test/internal/pkg/models"
	"github.com/Utro-tvar/medods-test/internal/tokens"
	"golang.org/x/crypto/bcrypt"
)

type Storage interface {
	Store(string) error
	CheckAndRemove(string) (exists bool, err error)
}

type TokenService struct {
	logger     *slog.Logger
	storage    Storage
	accessTTL  time.Duration
	refreshTTL time.Duration
	key        []byte
}

type TokensPair struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

func (t *TokensPair) ToJSON() (string, error) {
	res, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func New(logger *slog.Logger, storage Storage) *TokenService {
	return &TokenService{logger: logger, storage: storage}
}

func (t *TokenService) Generate(user models.User) TokensPair {
	const op = "service.Generate"
	access, refresh, err := tokens.Generate(user, t.accessTTL, t.refreshTTL, t.key)

	if err != nil {
		t.logger.Error("%s: Cannot generate tokens for user %s", op, user.GUID, slog.Any("error", err))
		return TokensPair{}
	}

	t.storage.Store(TokenHash(refresh))

	return TokensPair{Access: access, Refresh: refresh}
}

func (t *TokenService) Refresh(tokensJSON string) TokensPair {
	const op = "service.Refresh"
	tokensPair := &TokensPair{}
	err := json.Unmarshal([]byte(tokensJSON), tokensPair)
	if err != nil {
		t.logger.Error("%s: Cannot unmarshal tokens", op, slog.Any("error", err))
		return TokensPair{}
	}

	user, err := tokens.ExtractUser(tokensPair.Access, t.key)
	if err != nil {
		t.logger.Error("%s: Cannot extract user from token", op, slog.Any("error", err))
		return TokensPair{}
	}

	valid, err := tokens.Validate(tokensPair.Access, tokensPair.Refresh, t.key)
	if !valid {
		if errors.Is(err, tokens.ErrTokensPairIsInvalid) {
			t.logger.Info("%s: Tokens for user %s in invalid", op, user.GUID)
		} else {
			t.logger.Error("%s: Cannot validate tokens", op, slog.Any("error", err))
		}
		return TokensPair{}
	}

	hasToken, err := t.storage.CheckAndRemove(TokenHash(tokensPair.Refresh))
	if err != nil {
		t.logger.Error("%s: Error while talk to database", op, slog.Any("error", err))
	}
	if !hasToken {
		t.logger.Info("%s: Refresh token does not exist, user: %s", op, user.GUID)
		return TokensPair{}
	}
	return t.Generate(user)
}

func TokenHash(token string) string {
	bytes := []byte(token)
	bytes, _ = bcrypt.GenerateFromPassword(bytes[len(bytes)-70:], 12) // use 70 last bytes to generate hash
	return string(bytes)
}
