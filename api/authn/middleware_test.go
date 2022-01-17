package authn

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/classic-massok/classic-massok-be/api/authn/authnfakes"
	bizmodels "github.com/classic-massok/classic-massok-be/business/models"
	"github.com/classic-massok/classic-massok-be/lib"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestAuthnMW_ValidateToken_accessToken(t *testing.T) {
	authnMW := &AuthnMW{
		&authnfakes.FakeUserGetter{
			GetStub: func(c context.Context, s string) (*bizmodels.User, error) {
				return &bizmodels.User{
					ID: s,
					CusKeys: map[string]string{
						"1": "a",
						"2": "b",
					},
					Roles: bizmodels.Roles{
						"fake.role.one",
						"fake.role.two",
					},
				}, nil
			},
		},
	}

	userID := "user_id"

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	c.Set(lib.AccessTokenPrvKeyPathKey, "access_private_key.pem")

	at, _, err := GenerateAccessToken(c, userID)
	require.NoError(t, err)

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", at))
	res = httptest.NewRecorder()
	c = e.NewContext(req, res)

	c.Set(lib.AccessTokenPubKeyPathKey, "access_public_key.pem")

	handler := authnMW.ValidateToken(echo.HandlerFunc(func(c echo.Context) error {
		return nil
	}))

	require.NoError(t, handler(c))
	require.Equal(t, userID, c.Get(lib.UserIDKey).(string))
	require.Equal(t, AccessTokenType, c.Get(lib.TokenTypeKey).(string))
	require.Equal(t, map[string]string{"1": "a", "2": "b"}, c.Get(lib.CusKeysKey))
	require.Equal(t, bizmodels.Roles{"fake.role.one", "fake.role.two"}, c.Get(lib.RolesKey))
}

func TestAuthnMW_ValidateToken_refreshToken(t *testing.T) {
	authnMW := &AuthnMW{
		&authnfakes.FakeUserGetter{
			GetStub: func(c context.Context, s string) (*bizmodels.User, error) {
				return &bizmodels.User{
					ID: s,
					CusKeys: map[string]string{
						"1": "a",
						"2": "b",
					},
				}, nil
			},
		},
	}

	userID := "user_id"

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	c.Set(lib.RefreshTokenPrvKeyPathKey, "refresh_private_key.pem")

	rt, _, err := GenerateRefreshToken(c, userID, "a")
	require.NoError(t, err)

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", rt))
	res = httptest.NewRecorder()
	c = e.NewContext(req, res)

	c.Echo().IPExtractor = func(*http.Request) string { return "1" }
	c.Set(lib.AccessTokenPubKeyPathKey, "access_public_key.pem")
	c.Set(lib.RefreshTokenPubKeyPathKey, "refresh_public_key.pem")

	handler := authnMW.ValidateToken(echo.HandlerFunc(func(c echo.Context) error {
		return nil
	}))

	require.NoError(t, handler(c))
	require.Equal(t, userID, c.Get(lib.UserIDKey).(string))
	require.Equal(t, RefreshTokenType, c.Get(lib.TokenTypeKey).(string))
	require.Equal(t, map[string]string{"1": "a", "2": "b"}, c.Get(lib.CusKeysKey))
}

func TestAuthnMW_ValidateToken_accessToken_userNotFound(t *testing.T) {
	authnMW := &AuthnMW{
		&authnfakes.FakeUserGetter{
			GetStub: func(c context.Context, s string) (*bizmodels.User, error) {
				return nil, fmt.Errorf("user not found")
			},
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	c.Set(lib.AccessTokenPrvKeyPathKey, "access_private_key.pem")

	at, _, err := GenerateAccessToken(c, "user_id")
	require.NoError(t, err)

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", at))
	res = httptest.NewRecorder()
	c = e.NewContext(req, res)

	c.Set(lib.AccessTokenPubKeyPathKey, "access_public_key.pem")

	handler := authnMW.ValidateToken(echo.HandlerFunc(func(c echo.Context) error {
		return nil
	}))

	require.NoError(t, handler(c))
	require.NoError(t, handler(c))
	require.Equal(t, nil, c.Get(lib.UserIDKey))
	require.Equal(t, nil, c.Get(lib.TokenTypeKey))
	require.Equal(t, nil, c.Get(lib.CusKeysKey))
}

func TestAuthnMW_ValidateToken_accessToken_expired(t *testing.T) {
	authnMW := &AuthnMW{
		// User return required to ensure proper failure point
		&authnfakes.FakeUserGetter{
			GetStub: func(c context.Context, s string) (*bizmodels.User, error) {
				return &bizmodels.User{
					ID: s,
					CusKeys: map[string]string{
						"1": "a",
						"2": "b",
					},
				}, nil
			},
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	c.Set(lib.AccessTokenPrvKeyPathKey, "access_private_key.pem")

	at, _, err := GenerateAccessToken(c, "user_id")
	require.NoError(t, err)

	now := time.Now()
	// modify current time
	monkey.Patch(time.Now, func() time.Time {
		return now.Add(15 * time.Minute).Add(10 * time.Second)
	})
	defer monkey.Unpatch(time.Now)

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", at))
	res = httptest.NewRecorder()
	c = e.NewContext(req, res)

	c.Set(lib.AccessTokenPubKeyPathKey, "access_public_key.pem")

	handler := authnMW.ValidateToken(echo.HandlerFunc(func(c echo.Context) error {
		return nil
	}))

	require.NoError(t, handler(c))
	require.Equal(t, nil, c.Get(lib.UserIDKey))
	require.Equal(t, nil, c.Get(lib.TokenTypeKey))
	require.Equal(t, nil, c.Get(lib.CusKeysKey))
}

func TestAuthnMW_ValidateToken_noAuth(t *testing.T) {
	authnMW := &AuthnMW{
		&authnfakes.FakeUserGetter{
			// User return required to ensure proper failure point
			GetStub: func(c context.Context, s string) (*bizmodels.User, error) {
				return &bizmodels.User{
					ID: s,
					CusKeys: map[string]string{
						"1": "a",
						"2": "b",
					},
				}, nil
			},
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	c.Set(lib.AccessTokenPrvKeyPathKey, "access_private_key.pem")

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	res = httptest.NewRecorder()
	c = e.NewContext(req, res)

	c.Set(lib.AccessTokenPubKeyPathKey, "access_public_key.pem")

	handler := authnMW.ValidateToken(echo.HandlerFunc(func(c echo.Context) error {
		return nil
	}))

	require.NoError(t, handler(c))
	require.Equal(t, nil, c.Get(lib.UserIDKey))
	require.Equal(t, nil, c.Get(lib.TokenTypeKey))
	require.Equal(t, nil, c.Get(lib.CusKeysKey))
}

func TestAuthnMW_ValidateToken_badSignFile(t *testing.T) {
	authnMW := &AuthnMW{
		&authnfakes.FakeUserGetter{
			// User return required to ensure proper failure point
			GetStub: func(c context.Context, s string) (*bizmodels.User, error) {
				return &bizmodels.User{
					ID: s,
					CusKeys: map[string]string{
						"1": "a",
						"2": "b",
					},
				}, nil
			},
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	c.Set(lib.RefreshTokenPrvKeyPathKey, "refresh_private_key.pem")

	rt, _, err := GenerateRefreshToken(c, "user_id", "a")
	require.NoError(t, err)

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", rt))
	res = httptest.NewRecorder()
	c = e.NewContext(req, res)

	c.Set(lib.AccessTokenPubKeyPathKey, "")
	c.Set(lib.RefreshTokenPubKeyPathKey, "")

	handler := authnMW.ValidateToken(echo.HandlerFunc(func(c echo.Context) error {
		return nil
	}))

	require.NoError(t, handler(c))
	require.Equal(t, nil, c.Get(lib.UserIDKey))
	require.Equal(t, nil, c.Get(lib.TokenTypeKey))
	require.Equal(t, nil, c.Get(lib.CusKeysKey))
}

func TestAuthnMW_ValidateToken_refreshToken_badCusKey(t *testing.T) {
	authnMW := &AuthnMW{
		&authnfakes.FakeUserGetter{
			// User return required to ensure proper failure point
			GetStub: func(c context.Context, s string) (*bizmodels.User, error) {
				return &bizmodels.User{
					ID: s,
					CusKeys: map[string]string{
						"1": "a",
						"2": "b",
					},
				}, nil
			},
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	c.Set(lib.RefreshTokenPrvKeyPathKey, "refresh_private_key.pem")

	rt, _, err := GenerateRefreshToken(c, "user_id", "a")
	require.NoError(t, err)

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", rt))
	res = httptest.NewRecorder()
	c = e.NewContext(req, res)

	c.Echo().IPExtractor = func(*http.Request) string { return "0" }
	c.Set(lib.AccessTokenPubKeyPathKey, "access_public_key.pem")
	c.Set(lib.RefreshTokenPubKeyPathKey, "refresh_public_key.pem")

	handler := authnMW.ValidateToken(echo.HandlerFunc(func(c echo.Context) error {
		return nil
	}))

	require.NoError(t, handler(c))
	require.Equal(t, nil, c.Get(lib.UserIDKey))
	require.Equal(t, nil, c.Get(lib.TokenTypeKey))
	require.Equal(t, nil, c.Get(lib.CusKeysKey))
}

func TestAuthnMW_ValidateToken_refreshToken_userNotFound(t *testing.T) {
	authnMW := &AuthnMW{
		&authnfakes.FakeUserGetter{
			GetStub: func(c context.Context, s string) (*bizmodels.User, error) {
				return nil, fmt.Errorf("user not found")
			},
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	c.Set(lib.RefreshTokenPrvKeyPathKey, "refresh_private_key.pem")

	rt, _, err := GenerateRefreshToken(c, "user_id", "a")
	require.NoError(t, err)

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", rt))
	res = httptest.NewRecorder()
	c = e.NewContext(req, res)

	c.Echo().IPExtractor = func(*http.Request) string { return "1" }
	c.Set(lib.AccessTokenPubKeyPathKey, "access_public_key.pem")
	c.Set(lib.RefreshTokenPubKeyPathKey, "refresh_public_key.pem")

	handler := authnMW.ValidateToken(echo.HandlerFunc(func(c echo.Context) error {
		return nil
	}))

	require.NoError(t, handler(c))
	require.Equal(t, nil, c.Get(lib.UserIDKey))
	require.Equal(t, nil, c.Get(lib.TokenTypeKey))
	require.Equal(t, nil, c.Get(lib.CusKeysKey))
}

func TestAuthnMW_ValidateToken_refreshToken_expired(t *testing.T) {
	authnMW := &AuthnMW{
		&authnfakes.FakeUserGetter{
			// User return required to ensure proper failure point
			GetStub: func(c context.Context, s string) (*bizmodels.User, error) {
				return &bizmodels.User{
					ID: s,
					CusKeys: map[string]string{
						"1": "a",
						"2": "b",
					},
				}, nil
			},
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	c.Set(lib.RefreshTokenPrvKeyPathKey, "refresh_private_key.pem")

	rt, _, err := GenerateRefreshToken(c, "user_id", "a")
	require.NoError(t, err)

	now := time.Now()
	// modify current time
	monkey.Patch(time.Now, func() time.Time {
		return now.Add(18 * time.Hour).Add(10 * time.Second)
	})
	defer monkey.Unpatch(time.Now)

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", rt))
	res = httptest.NewRecorder()
	c = e.NewContext(req, res)

	c.Echo().IPExtractor = func(*http.Request) string { return "1" }
	c.Set(lib.AccessTokenPubKeyPathKey, "access_public_key.pem")
	c.Set(lib.RefreshTokenPubKeyPathKey, "refresh_public_key.pem")

	handler := authnMW.ValidateToken(echo.HandlerFunc(func(c echo.Context) error {
		return nil
	}))

	require.NoError(t, handler(c))
	require.Equal(t, nil, c.Get(lib.UserIDKey))
	require.Equal(t, nil, c.Get(lib.TokenTypeKey))
	require.Equal(t, nil, c.Get(lib.CusKeysKey))
}
