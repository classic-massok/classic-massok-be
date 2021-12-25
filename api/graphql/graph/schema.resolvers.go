package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/classic-massok/classic-massok-be/api/graphql/graph/generated"
	"github.com/classic-massok/classic-massok-be/api/graphql/graph/models"
	"github.com/classic-massok/classic-massok-be/business"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input models.CreateUserInput) (*models.CreateUserOutput, error) {
	usersBiz, err := business.NewUsersBiz()
	if err != nil {
		return nil, fmt.Errorf("error constructing users biz: %w", err)
	}

	id, err := usersBiz.New(ctx, "", input.Password, business.User{
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

func (r *mutationResolver) UpdateUser(ctx context.Context, input models.UpdateUserInput) (*models.User, error) {
	usersBiz, err := business.NewUsersBiz()
	if err != nil {
		return nil, fmt.Errorf("error constructing users biz: %w", err)
	}

	bizUser, err := usersBiz.Edit(ctx, input.ID, "", business.UserEdit{
		Email:     input.Email,
		Password:  input.Password,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Phone:     input.Phone,
		CanSMS:    input.CanSms,
		Birthday:  input.Birthday,
	})
	if err != nil {
		return nil, err
	}

	return &models.User{
		bizUser.GetID(),
		bizUser.Email,
		bizUser.FirstName,
		bizUser.LastName,
		bizUser.Phone,
		bizUser.CanSMS,
		bizUser.Birthday,
		bizUser.GetCreatedAt(),
		bizUser.GetUpdatedAt(),
		bizUser.GetCreatedBy(),
		bizUser.GetUpdatedBy(),
	}, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, input models.DeleteUserInput) (*models.DeleteUserOutput, error) {
	usersBiz, err := business.NewUsersBiz()
	if err != nil {
		return nil, fmt.Errorf("error constructing users biz: %w", err)
	}

	if err := usersBiz.Delete(ctx, input.ID, ""); err != nil {
		return nil, err
	}

	return &models.DeleteUserOutput{true}, nil
}

func (r *queryResolver) User(ctx context.Context, id string) (*models.User, error) {
	usersBiz, err := business.NewUsersBiz()
	if err != nil {
		return nil, fmt.Errorf("error constructing users biz: %w", err)
	}

	bizUser, err := usersBiz.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.User{
		bizUser.GetID(),
		bizUser.Email,
		bizUser.FirstName,
		bizUser.LastName,
		bizUser.Phone,
		bizUser.CanSMS,
		bizUser.Birthday,
		bizUser.GetCreatedAt(),
		bizUser.GetUpdatedAt(),
		bizUser.GetCreatedBy(),
		bizUser.GetUpdatedBy(),
	}, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*models.User, error) {
	usersBiz, err := business.NewUsersBiz()
	if err != nil {
		return nil, fmt.Errorf("error constructing users biz: %w", err)
	}

	bizUsers, err := usersBiz.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]*models.User, len(bizUsers))
	for i, bizUser := range bizUsers {
		users[i] = &models.User{
			bizUser.GetID(),
			bizUser.Email,
			bizUser.FirstName,
			bizUser.LastName,
			bizUser.Phone,
			bizUser.CanSMS,
			bizUser.Birthday,
			bizUser.GetCreatedAt(),
			bizUser.GetUpdatedAt(),
			bizUser.GetCreatedBy(),
			bizUser.GetUpdatedBy(),
		}
	}

	return users, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
