package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/classic-massok-be/api/graphql/graph/generated"
	"github.com/classic-massok-be/api/graphql/graph/models"
	"github.com/classic-massok-be/business"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input models.CreateUserInput) (*models.CreateUserOutput, error) {
	usersBiz, err := business.NewUsersBiz()
	if err != nil {
		return nil, fmt.Errorf("error constructing users biz: %w", err)
	}

	id, err := usersBiz.New(ctx, [12]byte{}, input.Password, business.User{
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Phone:     input.Phone,
		CanSMS:    input.CanSms,
		Birthday:  input.Birthday,
	})
	if err != nil {
		return nil, err
	}

	return &models.CreateUserOutput{
		ID: id,
	}, nil
}

func (r *queryResolver) User(ctx context.Context, id string) (*models.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Users(ctx context.Context) ([]*models.User, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
