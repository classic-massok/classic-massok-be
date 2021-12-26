package graphql

import (
	"context"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/classic-massok/classic-massok-be/api/graphql/generated"
	"github.com/classic-massok/classic-massok-be/api/graphql/resolvers"
	"github.com/labstack/echo/v4"
)

func Configure(graphql *echo.Group) {
	bindContext := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := context.WithValue(c.Request().Context(), struct{ Name string }{"echo"}, c)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}

	graphql.Use(bindContext)

	// Main graphql endpoiint
	graphql.POST("", graphqlMain)
	// Playground graphql endpoint
	graphql.GET("", graphqlPlayground)
}

func graphqlMain(c echo.Context) error {
	// acl := func(ctx context.Context, obj interface{}, next graph.Resolver, permission string) (res interface{}, err error) {
	// 	// Need to make a check permissions that uses the obj above to get the resource and it's perms
	// 	if err := authz.CheckPermissions(c, permission); err != nil {
	// 		return nil, err
	// 	}

	// 	return next(ctx)
	// }

	config := generated.Config{
		Resolvers:  &resolvers.Resolver{},
		Directives: generated.DirectiveRoot{},
		Complexity: generated.ComplexityRoot{},
	}

	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(config),
	)

	srv.ServeHTTP(c.Response(), c.Request())
	return nil
}

func graphqlPlayground(c echo.Context) error {
	srv := playground.Handler("GraphQL playground", "/api/graphql")
	srv.ServeHTTP(c.Response(), c.Request())
	return nil
}
