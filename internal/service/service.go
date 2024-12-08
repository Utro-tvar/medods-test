package service

import (
	"errors"
	"fmt"
	"log/slog"
	"net"

	"github.com/Utro-tvar/medods-test/internal/email"
	"github.com/Utro-tvar/medods-test/internal/pkg/models"
	"github.com/Utro-tvar/medods-test/internal/tokens"
	"golang.org/x/crypto/bcrypt"
)

type Storage interface {
	Store(string) error
	CheckAndRemove(string) (exists bool, err error)
}

type TokenService struct {
	logger  *slog.Logger
	storage Storage
	cfg     Config
}

func New(logger *slog.Logger, storage Storage) *TokenService {
	return &TokenService{logger: logger, storage: storage, cfg: ParseENV(logger)}
}

func (t *TokenService) Generate(user models.User) models.TokensPair {
	const op = "service.Generate"
	access, refresh, err := tokens.Generate(user, t.cfg.accessTTL, t.cfg.refreshTTL, t.cfg.key)

	if err != nil {
		t.logger.Error(fmt.Sprintf("%s: Cannot generate tokens for user %s", op, user.GUID), slog.Any("error", err))
		return models.TokensPair{}
	}

	t.storage.Store(TokenHash(refresh))

	return models.TokensPair{Access: access, Refresh: refresh}
}

func (t *TokenService) Refresh(tokensPair models.TokensPair, ip net.IP) models.TokensPair {
	const op = "service.Refresh"

	user, err := tokens.ExtractUser(tokensPair.Access, t.cfg.key)
	if err != nil {
		t.logger.Error(fmt.Sprintf("%s: Cannot extract user from token", op), slog.Any("error", err))
		return models.TokensPair{}
	}

	valid, err := tokens.Validate(tokensPair.Access, tokensPair.Refresh, t.cfg.key)
	if !valid {
		if errors.Is(err, tokens.ErrTokensPairIsInvalid) {
			t.logger.Info(fmt.Sprintf("%s: Tokens for user %s in invalid", op, user.GUID))
		} else {
			t.logger.Error(fmt.Sprintf("%s: Cannot validate tokens", op), slog.Any("error", err))
		}
		return models.TokensPair{}
	}

	hasToken, err := t.storage.CheckAndRemove(TokenHash(tokensPair.Refresh))
	if err != nil {
		t.logger.Error(fmt.Sprintf("%s: Error while talk to database", op), slog.Any("error", err))
	}
	if !hasToken {
		t.logger.Info(fmt.Sprintf("%s: Refresh token does not exist, user: %s", op, user.GUID))
		return models.TokensPair{}
	}
	if !net.IP.Equal(ip, user.IP) {
		email.Send("mock@email.com", []byte("Your IP has been changed"))
	}
	user.IP = ip
	return t.Generate(user)
}

func (t *TokenService) GetUser(access string) models.User {
	user, err := tokens.ExtractUser(access, t.cfg.key)
	if err != nil {
		t.logger.Error("cannot parse token", slog.Any("error", err))
	}
	return user
}

func TokenHash(token string) string {
	bytes := []byte(token)
	bytes, _ = bcrypt.GenerateFromPassword(bytes[len(bytes)-70:], 12) // use 70 last bytes to generate hash
	return string(bytes)
}
