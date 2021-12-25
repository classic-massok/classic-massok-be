package business

import (
	"context"
	"fmt"
	"time"

	"github.com/classic-massok-be/data/mongo"
	"github.com/classic-massok-be/data/mongo/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// NewUsersBiz is the constructure for the users business layers
func NewUsersBiz() (*usersBiz, error) {
	data, err := mongo.NewUsersData(utils.Database)
	if err != nil {
		return nil, fmt.Errorf("error constructing users biz: %w", err)
	}

	return &usersBiz{data}, nil
}

// usersBiz is the business layer for interacting with user data
type usersBiz struct {
	data usersData
}

//counterfeiter:generate . usersData
type usersData interface {
	New(ctx context.Context, loggedInUserID primitive.ObjectID, user mongo.User) (string, error)
}

type User struct {
	// Required fields
	Email     string
	FirstName string
	LastName  string

	// Optional fields
	Phone    *string
	CanSMS   *bool
	Birthday *time.Time
}

func (u *usersBiz) New(ctx context.Context, loggedInUserID primitive.ObjectID, password string, user User) (string, error) {
	passwordBytes := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return u.data.New(ctx, loggedInUserID, mongo.User{
		Email:     user.Email,
		Password:  hashedPassword,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		CanSMS:    user.CanSMS,
		Birthday:  user.Birthday,
	})
}
