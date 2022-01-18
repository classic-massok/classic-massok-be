package authz

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/classic-massok/classic-massok-be/api/authz/authzfakes"
	bizmodels "github.com/classic-massok/classic-massok-be/business/models"
	"github.com/classic-massok/classic-massok-be/lib"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestAuthzMW_LoadResource(t *testing.T) {
	id := "fake_resource"
	fakeResource := struct{ id string }{id}

	authzMW := &AuthzMW{
		&authzfakes.FakeAccessAllower{},
		&authzfakes.FakeResourceGetter{
			GetStub: func(c context.Context, s1, s2 string) (interface{}, error) {
				return struct{ id string }{s2}, nil
			},
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	handler := authzMW.LoadResource("", id)(echo.HandlerFunc(func(c echo.Context) error {
		return nil
	}))

	require.NoError(t, handler(c))
	require.Equal(t, c.Get(resourceKey).(struct{ id string }), fakeResource)
}

func TestAuthzMW_LoadResource_notFound(t *testing.T) {
	id := "fake_resource"

	authzMW := &AuthzMW{
		&authzfakes.FakeAccessAllower{},
		&authzfakes.FakeResourceGetter{
			GetStub: func(c context.Context, s1, s2 string) (interface{}, error) {
				return nil, fmt.Errorf("this is an error")
			},
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	handler := authzMW.LoadResource("", id)(echo.HandlerFunc(func(c echo.Context) error {
		return nil
	}))

	require.NoError(t, handler(c))
	require.Equal(t, c.Response().Status, 404)
	require.Nil(t, c.Get(resourceKey))
}

func TestAuthzMW_RequiresPermission(t *testing.T) {
	authzMW := &AuthzMW{
		&authzfakes.FakeAccessAllower{
			AccessAllowedStub: func(c context.Context, i interface{}, s1, s2 string, r bizmodels.Roles) (bool, error) {
				return true, nil
			},
		},
		&authzfakes.FakeResourceGetter{},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	c.Set(lib.RolesKey, bizmodels.Roles{"fake.role.one"})
	c.Set(lib.UserIDKey, "fake_user")

	handler := authzMW.RequiresPermission("fake-action")(echo.HandlerFunc(func(c echo.Context) error {
		return nil
	}))

	require.NoError(t, handler(c))
}

func TestAuthzMW_RequiresPermission_notAllowed(t *testing.T) {
	authzMW := &AuthzMW{
		&authzfakes.FakeAccessAllower{
			AccessAllowedStub: func(c context.Context, i interface{}, s1, s2 string, r bizmodels.Roles) (bool, error) {
				return false, nil
			},
		},
		&authzfakes.FakeResourceGetter{},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	c.Set(lib.RolesKey, bizmodels.Roles{"fake.role.one"})
	c.Set(lib.UserIDKey, "fake_user")

	handler := authzMW.RequiresPermission("fake-action")(echo.HandlerFunc(func(c echo.Context) error {
		return nil
	}))

	require.NoError(t, handler(c))
	require.Equal(t, c.Response().Status, 403)
}

func TestAuthzMW_RequiresPermission_accessError(t *testing.T) {
	authzMW := &AuthzMW{
		&authzfakes.FakeAccessAllower{
			AccessAllowedStub: func(c context.Context, i interface{}, s1, s2 string, r bizmodels.Roles) (bool, error) {
				return false, fmt.Errorf("some error")
			},
		},
		&authzfakes.FakeResourceGetter{},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	c.Set(lib.RolesKey, bizmodels.Roles{"fake.role.one"})
	c.Set(lib.UserIDKey, "fake_user")

	handler := authzMW.RequiresPermission("fake-action")(echo.HandlerFunc(func(c echo.Context) error {
		return nil
	}))

	require.NoError(t, handler(c))
	require.Equal(t, c.Response().Status, 500)
}

func TestAuthzMW_RequiresPermission_noRoles(t *testing.T) {
	authzMW := &AuthzMW{
		&authzfakes.FakeAccessAllower{
			AccessAllowedStub: func(c context.Context, i interface{}, s1, s2 string, r bizmodels.Roles) (bool, error) {
				return true, nil
			},
		},
		&authzfakes.FakeResourceGetter{},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	c.Set(lib.RolesKey, bizmodels.Roles{})
	c.Set(lib.UserIDKey, "fake_user")

	handler := authzMW.RequiresPermission("fake-action")(echo.HandlerFunc(func(c echo.Context) error {
		return nil
	}))

	require.NoError(t, handler(c))
	require.Equal(t, c.Response().Status, 403)
}

func TestAuthzMW_RequiresPermission_nilRoles(t *testing.T) {
	authzMW := &AuthzMW{
		&authzfakes.FakeAccessAllower{
			AccessAllowedStub: func(c context.Context, i interface{}, s1, s2 string, r bizmodels.Roles) (bool, error) {
				return true, nil
			},
		},
		&authzfakes.FakeResourceGetter{},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	c.Set(lib.UserIDKey, "fake_user")

	handler := authzMW.RequiresPermission("fake-action")(echo.HandlerFunc(func(c echo.Context) error {
		return nil
	}))

	require.NoError(t, handler(c))
	require.Equal(t, c.Response().Status, 401)
}

func TestAuthzMW_RequiresPermission_nilUser(t *testing.T) {
	authzMW := &AuthzMW{
		&authzfakes.FakeAccessAllower{
			AccessAllowedStub: func(c context.Context, i interface{}, s1, s2 string, r bizmodels.Roles) (bool, error) {
				return true, nil
			},
		},
		&authzfakes.FakeResourceGetter{},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	c.Set(lib.RolesKey, bizmodels.Roles{})

	handler := authzMW.RequiresPermission("fake-action")(echo.HandlerFunc(func(c echo.Context) error {
		return nil
	}))

	require.NoError(t, handler(c))
	require.Equal(t, c.Response().Status, 401)
}
