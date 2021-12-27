package resolvers

import (
	"context"
	"fmt"

	"github.com/classic-massok/classic-massok-be/api/authn"
	"github.com/classic-massok/classic-massok-be/api/graphql/models"
	"github.com/classic-massok/classic-massok-be/business"
)

func (m *mutation) CreateUser(ctx context.Context, input models.CreateUserInput) (*models.CreateUserOutput, error) {
	id, err := m.UsersBiz.New(ctx, "", input.Password, business.User{
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
	bizUser, err := q.UsersBiz.Get(ctx, id)
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
	c, err := echoContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if c.Get(authn.UserIDKey) == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	bizUsers, err := q.UsersBiz.GetAll(ctx)
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
	bizUser, err := m.UsersBiz.Edit(ctx, input.ID, "", false, business.UserEdit{
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
	if err := m.UsersBiz.Delete(ctx, input.ID, ""); err != nil {
		return nil, err
	}

	return &models.DeleteUserOutput{true}, nil
}
