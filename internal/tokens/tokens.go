package tokens

import (
	"errors"
	"math/rand"
	"net"
	"time"

	"github.com/Utro-tvar/medods-test/internal/pkg/models"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokensPairIsInvalid   = errors.New("tokens pair is invalid")
	ErrTokenPayloadIsInvalid = errors.New("token payload is invalid")
)

func Generate(userInfo models.User, accessTTL, refreshTTL time.Duration, key []byte) (access, refresh string, err error) {
	secret := GenerateSecret(10)
	jwtAccess := jwt.New(jwt.SigningMethodHS512)

	accessClaims := jwtAccess.Claims.(jwt.MapClaims)
	accessClaims["guid"] = userInfo.GUID
	accessClaims["ip"] = userInfo.IP.String()
	accessClaims["secret"] = secret
	accessClaims["exp"] = time.Now().Add(accessTTL).Unix()

	access, err = jwtAccess.SignedString(key)
	if err != nil {
		return "", "", err
	}

	jwtRefresh := jwt.New(jwt.SigningMethodHS512)

	refreshClaims := jwtRefresh.Claims.(jwt.MapClaims)
	refreshClaims["secret"] = secret
	refreshClaims["exp"] = time.Now().Add(refreshTTL).Unix()

	refresh, err = jwtRefresh.SignedString(key)
	if err != nil {
		return "", "", err
	}

	return access, refresh, err
}

func Validate(access, refresh string, key []byte) (bool, error) {
	keyFunk := func(*jwt.Token) (interface{}, error) { return key, nil }
	accessToken, err := jwt.Parse(access, keyFunk)
	if err != nil {
		return false, err
	}
	refreshToken, err := jwt.Parse(refresh, keyFunk)
	if err != nil {
		return false, err
	}

	accessClaims, ok := accessToken.Claims.(jwt.MapClaims)
	if !ok {
		return false, ErrTokenPayloadIsInvalid
	}
	refreshClaims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok {
		return false, ErrTokenPayloadIsInvalid
	}

	if refreshClaims["secret"] != accessClaims["secret"] {
		return false, ErrTokensPairIsInvalid
	}

	return true, nil
}

func ExtractUser(access string, key []byte) (models.User, error) {
	keyFunk := func(*jwt.Token) (interface{}, error) { return key, nil }
	token, err := jwt.Parse(access, keyFunk)
	if err != nil {
		return models.User{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return models.User{}, ErrTokenPayloadIsInvalid
	}
	ipstr, ok := claims["ip"].(string)
	if !ok {
		return models.User{}, ErrTokenPayloadIsInvalid
	}
	guid, ok := claims["guid"].(string)
	if !ok {
		return models.User{}, ErrTokenPayloadIsInvalid
	}
	ip := net.ParseIP(ipstr)
	return models.User{GUID: guid, IP: ip}, nil
}

func GenerateSecret(length int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	secret := make([]byte, length)

	for i := range secret {
		secret[i] = letterBytes[rand.Int63()%(int64)(len(letterBytes))]
	}
	return string(secret)
}
