package authn

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/classic-massok/classic-massok-be/api/authn/authnfakes"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestAuthnMW_ValidateToken(t *testing.T) {
	authnMW := &AuthnMW{
		&authnfakes.FakeUserGetter{},
	}

	// TODO: create fake token that contains user id returned in fake user getter
	// and check that validate token actually validates (happy paths and errors)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	handler := authnMW.ValidateToken(echo.HandlerFunc(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}))

	require.NoError(t, handler(c))
}
