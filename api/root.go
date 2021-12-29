package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/classic-massok/classic-massok-be/api/authn"
	"github.com/classic-massok/classic-massok-be/api/authz"
	"github.com/classic-massok/classic-massok-be/api/graphql"
	"github.com/classic-massok/classic-massok-be/api/rest"
	"github.com/classic-massok/classic-massok-be/business"
	"github.com/labstack/echo/v4"
)

type Router struct {
	ACLBiz   aclBiz
	UsersBiz usersBiz
}

func (r *Router) SetRouter() http.Handler {
	e := echo.New()
	authnMW := r.getAuthnMW()

	e.Use(authnMW.ValidateToken)
	e.Use(bindContext)

	apiRouter := e.Group("/api")

	restRouter := apiRouter.Group("/rest")
	rest.Configure(restRouter)

	graphqlRouter := apiRouter.Group("/graphql")
	r.getGraphQL().Configure(graphqlRouter)

	return e
}

func (r *Router) getAuthnMW() *authn.AuthnMW {
	return &authn.AuthnMW{
		r.UsersBiz,
	}
}

func (r *Router) getAuthzMW() *authz.AuthzMW {
	return &authz.AuthzMW{
		r.ACLBiz,
		&business.ResourceRepo{
			r.UsersBiz,
		},
	}
}

func (r *Router) getGraphQL() *graphql.GraphQL {
	return &graphql.GraphQL{
		r.UsersBiz,
	}
}

func bindContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := context.WithValue(c.Request().Context(), "EchoContext", c)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}

func EchoContextFromContext(ctx context.Context) (echo.Context, error) {
	echoContext := ctx.Value("EchoContext")
	if echoContext == nil {
		err := fmt.Errorf("could not retrieve echo.Context")
		return nil, err
	}

	ec, ok := echoContext.(echo.Context)
	if !ok {
		err := fmt.Errorf("echo.Context has wrong type")
		return nil, err
	}
	return ec, nil
}

type aclBiz interface {
	AccessAllowed(ctx context.Context, roles business.Roles, resource interface{}) (bool, error)
}

type usersBiz interface {
	Authn(ctx context.Context, email, password string) (string, string, error)
	New(ctx context.Context, loggedInUserID, password string, user business.User) (string, error)
	Get(ctx context.Context, id string) (*business.User, error)
	GetAll(ctx context.Context) ([]*business.User, error)
	Edit(ctx context.Context, id, loggedInUserID string, udpateCusKey bool, userEdit business.UserEdit) (*business.User, error)
	Delete(ctx context.Context, id, loggedInUserID string) error
}
