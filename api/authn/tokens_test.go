package authn

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/classic-massok/classic-massok-be/lib"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func Test_GenerateAccessToken(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	c := e.NewContext(req, res)
	c.Set(lib.AccessTokenPrvKeyPathKey, "access_private_key.pem")

	at, exp, err := GenerateAccessToken(c, "fake_user_id")
	assert.NotEmpty(t, at)
	assert.WithinDuration(t, time.Unix(exp, 0), time.Now().Add(15*time.Minute), 1*time.Second)
	assert.NoError(t, err)
}

func Test_GenerateRefreshToken(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	c := e.NewContext(req, res)
	c.Set(lib.RefreshTokenPrvKeyPathKey, "refresh_private_key.pem")

	at, exp, err := GenerateRefreshToken(c, "fake_user_id", "a")
	assert.NotEmpty(t, at)
	assert.WithinDuration(t, time.Unix(exp, 0), time.Now().Add(18*time.Hour), 1*time.Second)
	assert.NoError(t, err)
}

func Test_GenerateRefreshToken_bcryptError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	c := e.NewContext(req, res)
	c.Set(lib.RefreshTokenPrvKeyPathKey, "refresh_private_key.pem")

	err := fmt.Errorf("some error")
	monkey.Patch(bcrypt.GenerateFromPassword, func(_ []byte, _ int) ([]byte, error) {
		return nil, err
	})
	defer monkey.Unpatch(bcrypt.GenerateFromPassword)

	at, exp, err := GenerateRefreshToken(c, "fake_user_id", "a")
	assert.Empty(t, at)
	assert.Zero(t, exp)
	assert.EqualError(t, err, err.Error())
}

func Test_generateToken_badFilePath(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	c := e.NewContext(req, res)
	c.Set(lib.RefreshTokenPrvKeyPathKey, "bad_key.pem")

	_, _, err := GenerateRefreshToken(c, "fake_user_id", `¯\_(ツ)_/¯`)
	assert.EqualError(t, err, "error generating token: open bad_key.pem: no such file or directory")
}

func Test_generateToken_badSigningFile(t *testing.T) {
	ioutil.WriteFile("very_bad.pem", []byte("muy malo"), 0644)
	defer os.Remove("very_bad.pem")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	c := e.NewContext(req, res)
	c.Set(lib.RefreshTokenPrvKeyPathKey, "very_bad.pem")

	_, _, err := GenerateRefreshToken(c, "fake_user_id", `¯\_(ツ)_/¯`)
	assert.EqualError(t, err, "error generating token: Invalid Key: Key must be a PEM encoded PKCS1 or PKCS8 key")
}

func Test_validateToken_badSigningMethod(t *testing.T) {
	_, err := validateToken(&jwt.Token{}, "")
	require.EqualError(t, err, "Unexpected signing method in auth token")
}

func Test_validateToken_noSigningFile(t *testing.T) {
	_, err := validateToken(&jwt.Token{Method: &jwt.SigningMethodRSA{}}, "")
	require.EqualError(t, err, "open : no such file or directory")
}

func Test_validateToken_badSigningFile(t *testing.T) {
	ioutil.WriteFile("very_bad.pem", []byte("muy malo"), 0644)
	defer os.Remove("very_bad.pem")

	_, err := validateToken(&jwt.Token{Method: &jwt.SigningMethodRSA{}}, "very_bad.pem")
	require.EqualError(t, err, "Invalid Key: Key must be a PEM encoded PKCS1 or PKCS8 key")
}

func Test_extractToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Authorization", "Bearer Test_Token")
	require.NotEmpty(t, extractToken(req))
}

func Test_extractToken_noAuth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	require.Empty(t, extractToken(req))
}

func Test_extractToken_badAuth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Authorization", "Test_Token")
	require.Empty(t, extractToken(req))
}
