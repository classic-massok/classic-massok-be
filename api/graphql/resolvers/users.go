package resolvers

import (
	"context"

	"github.com/classic-massok/classic-massok-be/api/graphql/models"
	bizmodels "github.com/classic-massok/classic-massok-be/business/models"
)

func (m *mutation) CreateUser(ctx context.Context, input models.CreateUserInput) (*models.CreateUserOutput, error) {
	id, err := m.UsersBiz.New(ctx, "", input.Password, bizmodels.User{
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

	return &models.CreateUserOutput{
		ID: id,
	}, nil
}

func (q *query) User(ctx context.Context, input models.UserInput) (*models.User, error) {
	bizUser, err := q.UsersBiz.Get(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	panic("AH")

	return &models.User{
		bizUser.ID,
		bizUser.Email,
		bizUser.FirstName,
		bizUser.LastName,
		bizUser.Roles,
		bizUser.Phone,
		bizUser.CanSMS,
		bizUser.Birthday,
		bizUser.CreatedAt,
		bizUser.UpdatedAt,
		bizUser.CreatedBy,
		bizUser.UpdatedBy,
	}, nil
}

func (q *query) Users(ctx context.Context) ([]*models.User, error) {
	bizUsers, err := q.UsersBiz.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]*models.User, len(bizUsers))
	for i, bizUser := range bizUsers {
		users[i] = &models.User{
			bizUser.ID,
			bizUser.Email,
			bizUser.FirstName,
			bizUser.LastName,
			bizUser.Roles,
			bizUser.Phone,
			bizUser.CanSMS,
			bizUser.Birthday,
			bizUser.CreatedAt,
			bizUser.UpdatedAt,
			bizUser.CreatedBy,
			bizUser.UpdatedBy,
		}
	}

	return users, nil
}

func (m *mutation) UpdateUser(ctx context.Context, input models.UpdateUserInput) (*models.User, error) {
	bizUser, err := m.UsersBiz.Edit(ctx, input.ID, "", false, bizmodels.UserEdit{
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
		bizUser.ID,
		bizUser.Email,
		bizUser.FirstName,
		bizUser.LastName,
		bizUser.Roles,
		bizUser.Phone,
		bizUser.CanSMS,
		bizUser.Birthday,
		bizUser.CreatedAt,
		bizUser.UpdatedAt,
		bizUser.CreatedBy,
		bizUser.UpdatedBy,
	}, nil
}

func (m *mutation) DeleteUser(ctx context.Context, input models.DeleteUserInput) (*models.DeleteUserOutput, error) {
	if err := m.UsersBiz.Delete(ctx, input.ID, ""); err != nil {
		return nil, err
	}

	return &models.DeleteUserOutput{true}, nil
}
