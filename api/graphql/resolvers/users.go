package resolvers

import (
	"context"
	"fmt"

	"github.com/classic-massok/classic-massok-be/api/graphql/models"
	"github.com/classic-massok/classic-massok-be/business"
)

func (m *mutation) CreateUser(ctx context.Context, input models.CreateUserInput) (*models.CreateUserOutput, error) {
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

func (q *query) User(ctx context.Context, id string) (*models.User, error) {
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

func (q *query) Users(ctx context.Context) ([]*models.User, error) {
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

func (m *mutation) UpdateUser(ctx context.Context, input models.UpdateUserInput) (*models.User, error) {
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

func (m *mutation) DeleteUser(ctx context.Context, input models.DeleteUserInput) (*models.DeleteUserOutput, error) {
	usersBiz, err := business.NewUsersBiz()
	if err != nil {
		return nil, fmt.Errorf("error constructing users biz: %w", err)
	}

	if err := usersBiz.Delete(ctx, input.ID, ""); err != nil {
		return nil, err
	}

	return &models.DeleteUserOutput{true}, nil
}
