package resolvers

import (
	"context"
	"fmt"

	"github.com/classic-massok/classic-massok-be/api/graphql/generated"
	"github.com/labstack/echo/v4"
)

type Resolver struct{}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutation{}
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver {
	return &query{}
}

type mutation struct{}

type query struct{}

func echoContextFromContext(ctx context.Context) (echo.Context, error) {
	echoContext := ctx.Value(struct{ name string }{"echo"})
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
