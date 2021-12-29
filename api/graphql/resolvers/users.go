package resolvers

import (
	"context"

	graphqlmodels "github.com/classic-massok/classic-massok-be/api/graphql/models"
	"github.com/classic-massok/classic-massok-be/business"
)

func (m *mutation) CreateUser(ctx context.Context, input graphqlmodels.CreateUserInput) (*graphqlmodels.CreateUserOutput, error) {
	id, err := m.UsersBiz.New(ctx, "", input.Password, business.User{
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Roles:     input.Roles,
		Phone:     input.Phone,
		CanSMS:    input.CanSms,
		Birthday:  input.Birthday,
	})
	if err != nil {
		return nil, err
	}

	return &graphqlmodels.CreateUserOutput{
		ID: id,
	}, nil
}

func (q *query) User(ctx context.Context, input graphqlmodels.UserInput) (*graphqlmodels.User, error) {
	bizUser, err := q.UsersBiz.Get(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return &graphqlmodels.User{
		bizUser.GetID(),
		bizUser.Email,
		bizUser.FirstName,
		bizUser.LastName,
		bizUser.Roles,
		bizUser.Phone,
		bizUser.CanSMS,
		bizUser.Birthday,
		bizUser.GetCreatedAt(),
		bizUser.GetUpdatedAt(),
		bizUser.GetCreatedBy(),
		bizUser.GetUpdatedBy(),
	}, nil
}

func (q *query) Users(ctx context.Context) ([]*graphqlmodels.User, error) {
	bizUsers, err := q.UsersBiz.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]*graphqlmodels.User, len(bizUsers))
	for i, bizUser := range bizUsers {
		users[i] = &graphqlmodels.User{
			bizUser.GetID(),
			bizUser.Email,
			bizUser.FirstName,
			bizUser.LastName,
			bizUser.Roles,
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

func (m *mutation) UpdateUser(ctx context.Context, input graphqlmodels.UpdateUserInput) (*graphqlmodels.User, error) {
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

	return &graphqlmodels.User{
		bizUser.GetID(),
		bizUser.Email,
		bizUser.FirstName,
		bizUser.LastName,
		bizUser.Roles,
		bizUser.Phone,
		bizUser.CanSMS,
		bizUser.Birthday,
		bizUser.GetCreatedAt(),
		bizUser.GetUpdatedAt(),
		bizUser.GetCreatedBy(),
		bizUser.GetUpdatedBy(),
	}, nil
}

func (m *mutation) DeleteUser(ctx context.Context, input graphqlmodels.DeleteUserInput) (*graphqlmodels.DeleteUserOutput, error) {
	if err := m.UsersBiz.Delete(ctx, input.ID, ""); err != nil {
		return nil, err
	}

	return &graphqlmodels.DeleteUserOutput{true}, nil
}
