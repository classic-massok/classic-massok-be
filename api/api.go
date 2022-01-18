package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/classic-massok/classic-massok-be/api/authn"
	"github.com/classic-massok/classic-massok-be/api/authz"
	"github.com/classic-massok/classic-massok-be/api/graphql"
	"github.com/classic-massok/classic-massok-be/api/rest"
	bizmodels "github.com/classic-massok/classic-massok-be/business/models"
	"github.com/classic-massok/classic-massok-be/config"
	"github.com/classic-massok/classic-massok-be/lib"
	"github.com/labstack/echo/v4"
)

type Router struct {
	ACLBiz   aclBiz
	UsersBiz usersBiz
}

func (r *Router) SetRouter(resourceRepoBiz resourceRepoBiz, cfg *config.Config) http.Handler { // TODO: should this be "Set" or "Get"
	e := echo.New()
	e.IPExtractor = echo.ExtractIPDirect()

	authnMW := r.getAuthnMW()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(lib.AccessTokenPrvKeyPathKey, cfg.Tokens.AccessTokenPrvKeyPath)
			c.Set(lib.AccessTokenPubKeyPathKey, cfg.Tokens.AccessTokenPubKeyPath)
			c.Set(lib.RefreshTokenPrvKeyPathKey, cfg.Tokens.RefreshTokenPrvKeyPath)
			c.Set(lib.RefreshTokenPubKeyPathKey, cfg.Tokens.RefreshTokenPubKeyPath)
			ctx := context.WithValue(c.Request().Context(), lib.EchoContextKey, c)
			ctx = context.WithValue(ctx, lib.IPAddressKey, c.Echo().IPExtractor(c.Request()))
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	})
	e.Use(authnMW.ValidateToken)

	apiRouter := e.Group("/api")

	restRouter := apiRouter.Group("/rest")
	r.getRest(resourceRepoBiz).Configure(restRouter)

	graphqlRouter := apiRouter.Group("/graphql")
	r.getGraphQL(resourceRepoBiz, cfg).Configure(graphqlRouter)

	return e
}

func (r *Router) getAuthnMW() *authn.AuthnMW {
	return &authn.AuthnMW{
		r.UsersBiz,
	}
}

func (r *Router) getAuthzMW(resourceRepoBiz resourceRepoBiz) *authz.AuthzMW {
	return &authz.AuthzMW{
		r.ACLBiz,
		resourceRepoBiz,
	}
}

func (r *Router) getRest(resourceRepoBiz resourceRepoBiz) *rest.Rest {
	return &rest.Rest{
		r.getAuthzMW(resourceRepoBiz),
		r.UsersBiz,
	}
}

func (r *Router) getGraphQL(resourceRepoBiz resourceRepoBiz, cfg *config.Config) *graphql.GraphQL {
	return &graphql.GraphQL{
		r.ACLBiz,
		resourceRepoBiz,
		r.UsersBiz,
		cfg,
	}
}

func echoContextFromContext(ctx context.Context) (echo.Context, error) {
	echoContext := ctx.Value(lib.EchoContextKey)
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
	AccessAllowed(ctx context.Context, resource interface{}, action, userID string, roles bizmodels.Roles) (bool, error)
}

type resourceRepoBiz interface {
	Get(ctx context.Context, resourceType, resourceID string) (interface{}, error)
}

type usersBiz interface {
	Authn(ctx context.Context, email, password string) (string, map[string]string, error)
	New(ctx context.Context, loggedInUserID, password string, user bizmodels.User) (string, error)
	Get(ctx context.Context, id string) (*bizmodels.User, error)
	GetAll(ctx context.Context) ([]*bizmodels.User, error)
	Edit(ctx context.Context, id, loggedInUserID string, udpateCusKey bool, userEdit bizmodels.UserEdit) (*bizmodels.User, error)
	Delete(ctx context.Context, id, loggedInUserID string) error
}
