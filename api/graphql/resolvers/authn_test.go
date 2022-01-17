package resolvers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	graphqlmodels "github.com/classic-massok/classic-massok-be/api/graphql/models"
	"github.com/classic-massok/classic-massok-be/api/graphql/resolvers/resolversfakes"
	bizmodels "github.com/classic-massok/classic-massok-be/business/models"
	"github.com/classic-massok/classic-massok-be/lib"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestMutation_Login(t *testing.T) {
	userID := "user_id"
	cusKeys := map[string]string{
		"a": "1",
		"b": "2",
	}

	m := &mutation{
		&Resolver{
			&resolversfakes.FakeUsersBiz{
				AuthnStub: func(c context.Context, s1, s2 string) (string, map[string]string, error) {
					return userID, cusKeys, nil
				},
				EditStub: func(c context.Context, s1, s2 string, b bool, ue bizmodels.UserEdit) (*bizmodels.User, error) {
					return &bizmodels.User{
						CusKeys: cusKeys,
					}, nil
				},
			},
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	c.Set(lib.AccessTokenPrvKeyPathKey, "../../authn/access_private_key.pem")
	c.Set(lib.RefreshTokenPrvKeyPathKey, "../../authn/refresh_private_key.pem")

	ctx := context.WithValue(context.Background(), lib.IPAddressKey, "1")
	ctx = context.WithValue(ctx, lib.EchoContextKey, c)

	authOutput, err := m.Login(ctx, graphqlmodels.LoginInput{"", ""})

	require.NoError(t, err)
	require.NotNil(t, authOutput)
	require.NotEmpty(t, authOutput.AccessToken)
	require.NotEmpty(t, authOutput.RefreshToken)
	require.NotZero(t, authOutput.AccessTokenExpiry)
	require.NotZero(t, authOutput.RefreshTokenExpiry)
	require.Equal(t, c.Get(lib.UserIDKey).(string), userID)
}

func TestMutation_Login_badAuth(t *testing.T) {
	expectedErr := fmt.Errorf("some error")
	m := &mutation{
		&Resolver{
			&resolversfakes.FakeUsersBiz{
				AuthnStub: func(c context.Context, s1, s2 string) (string, map[string]string, error) {
					return "", nil, expectedErr
				},
			},
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	ctx := context.WithValue(context.Background(), lib.IPAddressKey, "1")
	ctx = context.WithValue(ctx, lib.EchoContextKey, c)

	authOutput, err := m.Login(ctx, graphqlmodels.LoginInput{"", ""})

	require.EqualError(t, err, expectedErr.Error())
	require.Nil(t, authOutput)
}

func TestMutation_Login_editError(t *testing.T) {
	userID := "user_id"
	cusKeys := map[string]string{
		"a": "1",
		"b": "2",
	}

	expectedErr := fmt.Errorf("some error")
	m := &mutation{
		&Resolver{
			&resolversfakes.FakeUsersBiz{
				AuthnStub: func(c context.Context, s1, s2 string) (string, map[string]string, error) {
					return userID, cusKeys, nil
				},
				EditStub: func(c context.Context, s1, s2 string, b bool, ue bizmodels.UserEdit) (*bizmodels.User, error) {
					return nil, expectedErr
				},
			},
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	ctx := context.WithValue(context.Background(), lib.IPAddressKey, "1")
	ctx = context.WithValue(ctx, lib.EchoContextKey, c)

	authOutput, err := m.Login(ctx, graphqlmodels.LoginInput{"", ""})

	require.ErrorAs(t, err, &expectedErr)
	require.Nil(t, authOutput)
}

func TestMutation_Login_accessTokenError(t *testing.T) {
	userID := "user_id"
	cusKeys := map[string]string{
		"a": "1",
		"b": "2",
	}

	m := &mutation{
		&Resolver{
			&resolversfakes.FakeUsersBiz{
				AuthnStub: func(c context.Context, s1, s2 string) (string, map[string]string, error) {
					return userID, cusKeys, nil
				},
				EditStub: func(c context.Context, s1, s2 string, b bool, ue bizmodels.UserEdit) (*bizmodels.User, error) {
					return &bizmodels.User{
						CusKeys: cusKeys,
					}, nil
				},
			},
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	c.Set(lib.AccessTokenPrvKeyPathKey, "")
	c.Set(lib.RefreshTokenPrvKeyPathKey, "../../authn/refresh_private_key.pem")

	ctx := context.WithValue(context.Background(), lib.IPAddressKey, "1")
	ctx = context.WithValue(ctx, lib.EchoContextKey, c)

	authOutput, err := m.Login(ctx, graphqlmodels.LoginInput{"", ""})

	require.EqualError(t, err, "error logging in: error generating token: open : no such file or directory")
	require.Nil(t, authOutput)
}

func TestMutation_Login_refreshTokenError(t *testing.T) {
	userID := "user_id"
	cusKeys := map[string]string{
		"a": "1",
		"b": "2",
	}

	m := &mutation{
		&Resolver{
			&resolversfakes.FakeUsersBiz{
				AuthnStub: func(c context.Context, s1, s2 string) (string, map[string]string, error) {
					return userID, cusKeys, nil
				},
				EditStub: func(c context.Context, s1, s2 string, b bool, ue bizmodels.UserEdit) (*bizmodels.User, error) {
					return &bizmodels.User{
						CusKeys: cusKeys,
					}, nil
				},
			},
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	c.Set(lib.AccessTokenPrvKeyPathKey, "../../authn/access_private_key.pem")
	c.Set(lib.RefreshTokenPrvKeyPathKey, "")

	ctx := context.WithValue(context.Background(), lib.IPAddressKey, "1")
	ctx = context.WithValue(ctx, lib.EchoContextKey, c)

	authOutput, err := m.Login(ctx, graphqlmodels.LoginInput{"", ""})

	require.EqualError(t, err, "error logging in: error generating token: open : no such file or directory")
	require.Nil(t, authOutput)
}
