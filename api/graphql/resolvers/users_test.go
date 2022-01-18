package resolvers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/classic-massok/classic-massok-be/api/graphql/models"
	"github.com/classic-massok/classic-massok-be/api/graphql/resolvers/resolversfakes"
	bizmodels "github.com/classic-massok/classic-massok-be/business/models"
	"github.com/stretchr/testify/require"
)

func TestMutation_CreateUser(t *testing.T) {
	userID := "user_id"
	m := &mutation{
		&Resolver{
			&resolversfakes.FakeUsersBiz{
				NewStub: func(c context.Context, s1, s2 string, u bizmodels.User) (string, error) {
					return userID, nil
				},
			},
		},
	}

	userOutput, err := m.CreateUser(context.Background(), models.CreateUserInput{})
	require.NoError(t, err)
	require.NotNil(t, userOutput)
	require.Equal(t, userOutput.ID, userID)
}

func TestMutation_CreateUser_NewError(t *testing.T) {
	expectedErr := fmt.Errorf("some error")
	m := &mutation{
		&Resolver{
			&resolversfakes.FakeUsersBiz{
				NewStub: func(c context.Context, s1, s2 string, u bizmodels.User) (string, error) {
					return "", expectedErr
				},
			},
		},
	}

	userOutput, err := m.CreateUser(context.Background(), models.CreateUserInput{})
	require.EqualError(t, err, expectedErr.Error())
	require.Nil(t, userOutput)
}

func TestQuery_User(t *testing.T) {
	now := time.Now()

	userID := "user_id"
	email := "fake@email.com"
	firstName := "Fake"
	lastName := "Name"
	roles := bizmodels.Roles{"fake.role.one", "fake.role.two"}
	phone := aws.String("8008675309")
	canSMS := aws.Bool(false)
	birthday := aws.Time(now.Add(69 * 24 * time.Hour))
	createdAt := now
	updatedAt := now.Add(5 * 24 * time.Hour)
	createdBy := "fake_user_id"
	updatedBy := "another_fake_user_id"

	q := &query{
		&Resolver{
			&resolversfakes.FakeUsersBiz{
				GetStub: func(c context.Context, s string) (*bizmodels.User, error) {
					return &bizmodels.User{
						userID, map[string]string{}, email, firstName, lastName,
						roles, phone, canSMS, birthday, bizmodels.Accounting{
							createdAt, updatedAt, createdBy, updatedBy,
						},
					}, nil
				},
			},
		},
	}

	userOutput, err := q.User(context.Background(), models.UserInput{})
	require.NoError(t, err)
	require.NotNil(t, userOutput)
	require.EqualValues(t, userOutput, &models.User{
		userID, email, firstName, lastName, roles, phone, canSMS, birthday, createdAt,
		updatedAt, createdBy, updatedBy,
	})
}

func TestQuery_User_GetError(t *testing.T) {
	expectedErr := fmt.Errorf("some error")

	q := &query{
		&Resolver{
			&resolversfakes.FakeUsersBiz{
				GetStub: func(c context.Context, s string) (*bizmodels.User, error) {
					return nil, expectedErr
				},
			},
		},
	}

	userOutput, err := q.User(context.Background(), models.UserInput{})
	require.EqualError(t, err, expectedErr.Error())
	require.Nil(t, userOutput)
}

func TestQuery_Users(t *testing.T) {
	now := time.Now()

	userID := "user_id"
	email := "fake@email.com"
	firstName := "Fake"
	lastName := "Name"
	roles := bizmodels.Roles{"fake.role.one", "fake.role.two"}
	phone := aws.String("8008675309")
	canSMS := aws.Bool(false)
	birthday := aws.Time(now.Add(69 * 24 * time.Hour))
	createdAt := now
	updatedAt := now.Add(5 * 24 * time.Hour)
	createdBy := "fake_user_id"
	updatedBy := "another_fake_user_id"

	q := &query{
		&Resolver{
			&resolversfakes.FakeUsersBiz{
				GetAllStub: func(c context.Context) ([]*bizmodels.User, error) {
					return []*bizmodels.User{{
						userID, map[string]string{}, email, firstName, lastName,
						roles, phone, canSMS, birthday, bizmodels.Accounting{
							createdAt, updatedAt, createdBy, updatedBy,
						},
					}}, nil
				},
			},
		},
	}

	usersOutput, err := q.Users(context.Background())
	require.NoError(t, err)
	require.NotNil(t, usersOutput)
	require.EqualValues(t, usersOutput, []*models.User{{
		userID, email, firstName, lastName, roles, phone, canSMS, birthday, createdAt,
		updatedAt, createdBy, updatedBy,
	}})
}

func TestQuery_Users_GetAllError(t *testing.T) {
	expectedErr := fmt.Errorf("some error")

	q := &query{
		&Resolver{
			&resolversfakes.FakeUsersBiz{
				GetAllStub: func(c context.Context) ([]*bizmodels.User, error) {
					return nil, expectedErr
				},
			},
		},
	}

	usersOutput, err := q.Users(context.Background())
	require.EqualError(t, err, expectedErr.Error())
	require.Nil(t, usersOutput)
}

func TestMutation_UpdateUser(t *testing.T) {
	now := time.Now()

	userID := "user_id"
	email := "fake@email.com"
	firstName := "Fake"
	lastName := "Name"
	roles := bizmodels.Roles{"fake.role.one", "fake.role.two"}
	phone := aws.String("8008675309")
	canSMS := aws.Bool(false)
	birthday := aws.Time(now.Add(69 * 24 * time.Hour))
	createdAt := now
	updatedAt := now.Add(5 * 24 * time.Hour)
	createdBy := "fake_user_id"
	updatedBy := "another_fake_user_id"

	m := &mutation{
		&Resolver{
			&resolversfakes.FakeUsersBiz{
				EditStub: func(c context.Context, s1, s2 string, b bool, ue bizmodels.UserEdit) (*bizmodels.User, error) {
					return &bizmodels.User{
						userID, map[string]string{}, email, firstName, lastName,
						roles, phone, canSMS, birthday, bizmodels.Accounting{
							createdAt, updatedAt, createdBy, updatedBy,
						},
					}, nil
				},
			},
		},
	}

	userOutput, err := m.UpdateUser(context.Background(), models.UpdateUserInput{})
	require.NoError(t, err)
	require.NotNil(t, userOutput)
	require.EqualValues(t, userOutput, &models.User{
		userID, email, firstName, lastName, roles, phone, canSMS, birthday, createdAt,
		updatedAt, createdBy, updatedBy,
	})
}

func TestMutation_UpdateUser_EditError(t *testing.T) {
	expectedErr := fmt.Errorf("some error")

	m := &mutation{
		&Resolver{
			&resolversfakes.FakeUsersBiz{
				EditStub: func(c context.Context, s1, s2 string, b bool, ue bizmodels.UserEdit) (*bizmodels.User, error) {
					return nil, expectedErr
				},
			},
		},
	}

	userOutput, err := m.UpdateUser(context.Background(), models.UpdateUserInput{})
	require.EqualError(t, err, expectedErr.Error())
	require.Nil(t, userOutput)
}

func TestMutation_DeleteUser(t *testing.T) {
	m := &mutation{
		&Resolver{
			&resolversfakes.FakeUsersBiz{
				DeleteStub: func(c context.Context, s1, s2 string) error {
					return nil
				},
			},
		},
	}

	userOutput, err := m.DeleteUser(context.Background(), models.DeleteUserInput{})
	require.NoError(t, err)
	require.NotNil(t, userOutput)
	require.True(t, userOutput.Success)
}

func TestMutation_DeleteUser_DeleteError(t *testing.T) {
	expectedErr := fmt.Errorf("some error")

	m := &mutation{
		&Resolver{
			&resolversfakes.FakeUsersBiz{
				DeleteStub: func(c context.Context, s1, s2 string) error {
					return expectedErr
				},
			},
		},
	}

	userOutput, err := m.DeleteUser(context.Background(), models.DeleteUserInput{})
	require.EqualError(t, err, expectedErr.Error())
	require.Nil(t, userOutput)
}
