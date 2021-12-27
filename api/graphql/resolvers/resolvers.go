package resolvers

import (
	"context"
	"fmt"

	"github.com/classic-massok/classic-massok-be/api/graphql/generated"
	"github.com/classic-massok/classic-massok-be/business"
	"github.com/labstack/echo/v4"
)

type Resolver struct {
	UsersBiz usersBiz
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutation{r}
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver {
	return &query{r}
}

type mutation struct{ *Resolver }

type query struct{ *Resolver }

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

type usersBiz interface {
	Authn(ctx context.Context, email, password string) (string, string, error)
	New(ctx context.Context, loggedInUserID, password string, user business.User) (string, error)
	Get(ctx context.Context, id string) (*business.User, error)
	GetAll(ctx context.Context) ([]*business.User, error)
	Edit(ctx context.Context, id, loggedInUserID string, userEdit business.UserEdit) (*business.User, error)
	Delete(ctx context.Context, id, loggedInUserID string) error
}
