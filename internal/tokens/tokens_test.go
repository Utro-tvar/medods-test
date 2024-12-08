package tokens_test

import (
	"net"
	"testing"
	"time"

	"github.com/Utro-tvar/medods-test/internal/pkg/models"
	"github.com/Utro-tvar/medods-test/internal/tokens"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestGeneration(t *testing.T) {
	var (
		key     = []byte("askjda;ksjd;a;su")
		keyFunc = func(t *jwt.Token) (interface{}, error) { return key, nil }
		user    = models.User{GUID: "0", IP: net.ParseIP("192.168.0.1")}
	)
	t.Run("Regular", func(t *testing.T) {

		access, refresh, err := tokens.Generate(user, 10*time.Minute, 1000*time.Hour, key)

		require.Nil(t, err, "Error when generate tokens")

		accessToken, err := jwt.Parse(access, keyFunc)
		require.Nil(t, err, "Error while parsing access token")

		accessBody, ok := accessToken.Claims.(jwt.MapClaims)
		require.True(t, ok, "Access token body in wrong format")

		require.Equalf(t, user.GUID, accessBody["guid"], "Wrong guid: %s expected: %s", accessBody["guid"], user.GUID)
		ipstr, ok := accessBody["ip"].(string)
		require.True(t, ok, "IP addres encoded wrong")
		ip := net.ParseIP(ipstr)
		require.NotNilf(t, ip, "IP in wrong format: %s", ipstr)
		require.Equalf(t, user.IP, ip, "Wrong ip: %s expected: %s", ip.String(), user.IP.String())

		refreshToken, err := jwt.Parse(refresh, keyFunc)
		require.Nil(t, err, "Error while parsing refresh token")

		refreshBody, ok := refreshToken.Claims.(jwt.MapClaims)
		require.True(t, ok, "Refresh token body in wrong format")

		require.Equal(t, accessBody["secret"], refreshBody["secret"], "Tokens must have same secrets")
	})
	t.Run("Expired", func(t *testing.T) {
		access, refresh, err := tokens.Generate(user, -10*time.Minute, -1000*time.Hour, key)
		require.Nil(t, err, "Error when generate tokens")

		_, err = jwt.Parse(access, keyFunc)
		require.ErrorIs(t, err, jwt.ErrTokenExpired, "Access token must be expired")

		_, err = jwt.Parse(refresh, keyFunc)
		require.ErrorIs(t, err, jwt.ErrTokenExpired, "Refresh token must be expired")
	})
}
