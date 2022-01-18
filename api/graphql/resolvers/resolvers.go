package resolvers

import (
	"context"

	"github.com/classic-massok/classic-massok-be/api/graphql/generated"
	bizmodels "github.com/classic-massok/classic-massok-be/business/models"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

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

//counterfeiter:generate . usersBiz
type usersBiz interface {
	Authn(ctx context.Context, email, password string) (string, map[string]string, error)
	New(ctx context.Context, loggedInUserID, password string, user bizmodels.User) (string, error)
	Get(ctx context.Context, id string) (*bizmodels.User, error)
	GetAll(ctx context.Context) ([]*bizmodels.User, error)
	Edit(ctx context.Context, id, loggedInUserID string, updateCusKey bool, userEdit bizmodels.UserEdit) (*bizmodels.User, error)
	Delete(ctx context.Context, id, loggedInUserID string) error
}
