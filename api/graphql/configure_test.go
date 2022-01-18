package graphql

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/classic-massok/classic-massok-be/api/graphql/graphqlfakes"
	"github.com/classic-massok/classic-massok-be/business/models"
	bizmodels "github.com/classic-massok/classic-massok-be/business/models"
	"github.com/classic-massok/classic-massok-be/config"
	"github.com/classic-massok/classic-massok-be/lib"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestGraphQL_Configure(t *testing.T) {
	g := &GraphQL{}
	g.Configure(echo.New().Group("")) // TODO: is there more we need to assert here?
}

func TestGraphQL_graphqlMain(t *testing.T) {
	g := &GraphQL{}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	require.NoError(t, g.graphqlMain(c)) // TODO: is there more we need to assert here?
}

func TestGraphQL_graphqlPlayground(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	require.NoError(t, graphqlPlayground(c)) // TODO: is there more we need to assert here?
}

func TestGraphQL_loadResource(t *testing.T) {
	g := &GraphQL{
		&graphqlfakes.FakeAccessAllower{},
		&graphqlfakes.FakeResourceRepoBiz{
			GetStub: func(c context.Context, s1, s2 string) (interface{}, error) {
				return nil, nil
			},
		},
		&graphqlfakes.FakeUsersBiz{},
		&config.Config{},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	ctx := context.WithValue(c.Request().Context(), lib.EchoContextKey, c)
	nextFn := func(context.Context) (interface{}, error) { return "", nil }

	next, err := g.loadResource(ctx, map[string]interface{}{"id": ""}, nextFn, "")
	require.NoError(t, err)
	require.NotNil(t, next)
}

func TestGraphQL_loadResource_badObj(t *testing.T) {
	g := &GraphQL{
		&graphqlfakes.FakeAccessAllower{},
		&graphqlfakes.FakeResourceRepoBiz{
			GetStub: func(c context.Context, s1, s2 string) (interface{}, error) {
				return nil, nil
			},
		},
		&graphqlfakes.FakeUsersBiz{},
		&config.Config{},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	ctx := context.WithValue(c.Request().Context(), lib.EchoContextKey, c)
	nextFn := func(context.Context) (interface{}, error) { return nil, nil }

	next, err := g.loadResource(ctx, nil, nextFn, "")
	require.Errorf(t, err, "invalid  object for id getter")
	require.Nil(t, next)
}

func TestGraphQL_loadResource_LoadResourceError(t *testing.T) {
	g := &GraphQL{
		&graphqlfakes.FakeAccessAllower{},
		&graphqlfakes.FakeResourceRepoBiz{
			GetStub: func(c context.Context, s1, s2 string) (interface{}, error) {
				return nil, fmt.Errorf("some error")
			},
		},
		&graphqlfakes.FakeUsersBiz{},
		&config.Config{},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	ctx := context.WithValue(c.Request().Context(), lib.EchoContextKey, c)
	nextFn := func(context.Context) (interface{}, error) { return "", nil }

	next, err := g.loadResource(ctx, map[string]interface{}{"id": ""}, nextFn, "")
	require.EqualError(t, err, lib.ErrNotFound.Error())
	require.Nil(t, next)
}

func TestGraphQL_acl(t *testing.T) {
	g := &GraphQL{
		&graphqlfakes.FakeAccessAllower{
			AccessAllowedStub: func(c context.Context, i interface{}, s1, s2 string, r models.Roles) (bool, error) {
				return true, nil
			},
		},
		&graphqlfakes.FakeResourceRepoBiz{},
		&graphqlfakes.FakeUsersBiz{},
		&config.Config{},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	c.Set(lib.RolesKey, bizmodels.Roles{"fake.role.one"})
	c.Set(lib.UserIDKey, "")

	ctx := context.WithValue(c.Request().Context(), lib.EchoContextKey, c)
	nextFn := func(context.Context) (interface{}, error) { return "", nil }

	next, err := g.acl(ctx, nil, nextFn, "")
	require.NoError(t, err)
	require.NotNil(t, next)
}

func TestGraphQL_acl_noPermission(t *testing.T) {
	g := &GraphQL{
		&graphqlfakes.FakeAccessAllower{
			AccessAllowedStub: func(c context.Context, i interface{}, s1, s2 string, r models.Roles) (bool, error) {
				return false, nil
			},
		},
		&graphqlfakes.FakeResourceRepoBiz{},
		&graphqlfakes.FakeUsersBiz{},
		&config.Config{},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	c.Set(lib.RolesKey, bizmodels.Roles{"fake.role.one"})
	c.Set(lib.UserIDKey, "")

	ctx := context.WithValue(c.Request().Context(), lib.EchoContextKey, c)
	nextFn := func(context.Context) (interface{}, error) { return "", nil }

	next, err := g.acl(ctx, nil, nextFn, "")
	require.EqualError(t, err, lib.ErrForbidden.Error())
	require.Nil(t, next)
}

func TestGraphQL_recover_zeroValues(t *testing.T) {
	g := &GraphQL{
		&graphqlfakes.FakeAccessAllower{},
		&graphqlfakes.FakeResourceRepoBiz{},
		&graphqlfakes.FakeUsersBiz{},
		&config.Config{},
	}

	expectedErr := fmt.Errorf("some error")
	output, err := testCaptureOutput(func() error {
		return g.recover(context.Background(), expectedErr)
	})
	require.EqualError(t, err, lib.ErrServerError.Error())
	require.Empty(t, output)
}

func TestGraphQL_recover_StdOutPanics(t *testing.T) {
	g := &GraphQL{
		&graphqlfakes.FakeAccessAllower{},
		&graphqlfakes.FakeResourceRepoBiz{},
		&graphqlfakes.FakeUsersBiz{},
		&config.Config{
			Logging: struct {
				StdOutPanics bool "yaml:\"stdOutPanics\""
				HTTPVerbose  bool "yaml:\"httpVerbose\""
			}{true, false},
		},
	}

	expectedErr := fmt.Errorf("some error")
	output, err := testCaptureOutput(func() error {
		return g.recover(context.Background(), expectedErr)
	})
	require.EqualError(t, err, lib.ErrServerError.Error())
	require.NotEmpty(t, output) // TODO: find a way to require equal without worrying about addresses
}

func TestGraphQL_recover_HTTPVerbose(t *testing.T) {
	g := &GraphQL{
		&graphqlfakes.FakeAccessAllower{},
		&graphqlfakes.FakeResourceRepoBiz{},
		&graphqlfakes.FakeUsersBiz{},
		&config.Config{
			Logging: struct {
				StdOutPanics bool "yaml:\"stdOutPanics\""
				HTTPVerbose  bool "yaml:\"httpVerbose\""
			}{false, true},
		},
	}

	expectedErr := fmt.Errorf("some error")
	output, err := testCaptureOutput(func() error {
		return g.recover(context.Background(), expectedErr)
	})
	require.ErrorAs(t, err, &lib.ErrServerError)
	require.ErrorAs(t, err, &expectedErr)
	require.Empty(t, output)
}

func testCaptureOutput(fn func() error) (string, error) {
	var buf bytes.Buffer
	errOut = &buf
	err := fn()
	return buf.String(), err
}
