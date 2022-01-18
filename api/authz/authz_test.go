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

func Test_LoadResource(t *testing.T) {
	fakeResource := "fake_resource"
	bizRepo := &authzfakes.FakeResourceGetter{
		GetStub: func(c context.Context, s1, s2 string) (interface{}, error) {
			return fakeResource, nil
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	require.NoError(t, LoadResource(c, bizRepo, "", ""))
	require.Equal(t, fakeResource, c.Get(resourceKey).(string))
}

func Test_LoadResource_notFound(t *testing.T) {
	bizRepo := &authzfakes.FakeResourceGetter{
		GetStub: func(c context.Context, s1, s2 string) (interface{}, error) {
			return nil, lib.ErrNotFound
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	require.EqualError(t, LoadResource(c, bizRepo, "", ""), lib.ErrNotFound.Error())
	require.Equal(t, nil, c.Get(resourceKey))
}

func Test_RequiresPermission(t *testing.T) {
	aclBiz := &authzfakes.FakeAccessAllower{
		AccessAllowedStub: func(c context.Context, i interface{}, s1, s2 string, r bizmodels.Roles) (bool, error) {
			return true, nil
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	c.Set(lib.RolesKey, bizmodels.Roles{"fake.role.one"})
	c.Set(lib.UserIDKey, "fake_user")

	require.NoError(t, RequiresPermission(c, aclBiz, ""))
}

func Test_RequiresPermission_notAllowed(t *testing.T) {
	aclBiz := &authzfakes.FakeAccessAllower{
		AccessAllowedStub: func(c context.Context, i interface{}, s1, s2 string, r bizmodels.Roles) (bool, error) {
			return false, nil
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	c.Set(lib.RolesKey, bizmodels.Roles{"fake.role.one"})
	c.Set(lib.UserIDKey, "fake_user")

	require.EqualError(t, RequiresPermission(c, aclBiz, ""), lib.ErrForbidden.Error())
}

func Test_RequiresPermission_accessError(t *testing.T) {
	aclBiz := &authzfakes.FakeAccessAllower{
		AccessAllowedStub: func(c context.Context, i interface{}, s1, s2 string, r bizmodels.Roles) (bool, error) {
			return false, fmt.Errorf("some error")
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	c.Set(lib.RolesKey, bizmodels.Roles{"fake.role.one"})
	c.Set(lib.UserIDKey, "fake_user")

	require.ErrorAs(t, RequiresPermission(c, aclBiz, ""), &lib.ErrServerError)
}

func Test_RequiresPermission_noRoles(t *testing.T) {
	aclBiz := &authzfakes.FakeAccessAllower{
		AccessAllowedStub: func(c context.Context, i interface{}, s1, s2 string, r bizmodels.Roles) (bool, error) {
			return true, nil
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	c.Set(lib.RolesKey, bizmodels.Roles{})
	c.Set(lib.UserIDKey, "fake_user")

	require.EqualError(t, RequiresPermission(c, aclBiz, ""), lib.ErrForbidden.Error())
}

func Test_RequiresPermission_nilRoles(t *testing.T) {
	aclBiz := &authzfakes.FakeAccessAllower{
		AccessAllowedStub: func(c context.Context, i interface{}, s1, s2 string, r bizmodels.Roles) (bool, error) {
			return true, nil
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	c.Set(lib.UserIDKey, "fake_user")

	require.EqualError(t, RequiresPermission(c, aclBiz, ""), lib.ErrUnauthorized.Error())
}

func Test_RequiresPermission_nilUser(t *testing.T) {
	aclBiz := &authzfakes.FakeAccessAllower{
		AccessAllowedStub: func(c context.Context, i interface{}, s1, s2 string, r bizmodels.Roles) (bool, error) {
			return true, nil
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	c.Set(lib.RolesKey, bizmodels.Roles{"fake.role.one"})

	require.EqualError(t, RequiresPermission(c, aclBiz, ""), lib.ErrUnauthorized.Error())
}
