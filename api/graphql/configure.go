package graphql

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime/debug"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/classic-massok/classic-massok-be/api/authz"
	"github.com/classic-massok/classic-massok-be/api/graphql/generated"
	"github.com/classic-massok/classic-massok-be/api/graphql/resolvers"
	bizmodels "github.com/classic-massok/classic-massok-be/business/models"
	"github.com/classic-massok/classic-massok-be/config"
	"github.com/classic-massok/classic-massok-be/lib"
	"github.com/labstack/echo/v4"
)

//go:generate go run github.com/99designs/gqlgen --verbose
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

var errOut io.Writer = os.Stderr

type GraphQL struct {
	ACLBiz          accessAllower
	ResourceRepoBiz resourceRepoBiz
	UsersBiz        usersBiz

	Cfg *config.Config
}

func (g *GraphQL) Configure(graphql *echo.Group) {
	// Main graphql endpoint
	graphql.POST("", g.graphqlMain)
	// Playground graphql endpoint
	graphql.GET("", graphqlPlayground)
}

func (g *GraphQL) graphqlMain(c echo.Context) error {
	config := generated.Config{
		Resolvers:  g.buildResolver(),
		Directives: generated.DirectiveRoot{g.acl, g.loadResource},
		Complexity: generated.ComplexityRoot{},
	}

	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(config),
	)

	srv.SetRecoverFunc(g.recover)

	srv.ServeHTTP(c.Response(), c.Request())
	return nil
}

func (g *GraphQL) buildResolver() *resolvers.Resolver {
	return &resolvers.Resolver{
		g.UsersBiz,
	}
}

func (g *GraphQL) loadResource(ctx context.Context, obj interface{}, next graphql.Resolver, resourceType string) (interface{}, error) {
	c := ctx.Value(lib.EchoContextKey).(echo.Context)
	input, ok := obj.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid object for id getter: %T", obj) // TODO: figure out best error here (and a pattern for gql)
	}

	if err := authz.LoadResource(c, g.ResourceRepoBiz, resourceType, input["id"].(string)); err != nil {
		return nil, fmt.Errorf("not found")
	}

	return next(ctx)
}

func (g *GraphQL) acl(ctx context.Context, obj interface{}, next graphql.Resolver, action string) (interface{}, error) {
	c := ctx.Value(lib.EchoContextKey).(echo.Context)
	if err := authz.RequiresPermission(c, g.ACLBiz, action); err != nil {
		return nil, err
	}

	return next(ctx)
}

func (g *GraphQL) recover(ctx context.Context, err interface{}) error {
	if g.Cfg.Logging.StdOutPanics {
		fmt.Fprintln(errOut, err)
		fmt.Fprintln(errOut)
		fmt.Fprintln(errOut, string(debug.Stack()))
	}

	if g.Cfg.Logging.HTTPVerbose {
		return fmt.Errorf("%w: %v\n\n%s", lib.ErrServerError, err, string(debug.Stack()))
	}

	return lib.ErrServerError
}

func graphqlPlayground(c echo.Context) error {
	srv := playground.Handler("GraphQL playground", "/api/graphql")
	srv.ServeHTTP(c.Response(), c.Request())
	return nil
}

//counterfeiter:generate . accessAllower
type accessAllower interface {
	AccessAllowed(ctx context.Context, resource interface{}, action, userID string, roles bizmodels.Roles) (bool, error)
}

//counterfeiter:generate . resourceRepoBiz
type resourceRepoBiz interface {
	Get(ctx context.Context, resourceType, resourceID string) (interface{}, error)
}

//counterfeiter:generate . usersBiz
type usersBiz interface {
	Authn(ctx context.Context, email, password string) (string, map[string]string, error)
	New(ctx context.Context, loggedInUserID, password string, user bizmodels.User) (string, error)
	Get(ctx context.Context, id string) (*bizmodels.User, error)
	GetAll(ctx context.Context) ([]*bizmodels.User, error)
	Edit(ctx context.Context, id, loggedInUserID string, updateCusKey bool, userEdit bizmodels.UserEdit) (*bizmodels.User, error)
	Delete(ctx context.Context, id, loggedInUserID string) error
}
