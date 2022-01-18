// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

import (
	"time"
)

type AuthOutput struct {
	AccessToken        string `json:"accessToken"`
	AccessTokenExpiry  int64  `json:"accessTokenExpiry"`
	RefreshToken       string `json:"refreshToken"`
	RefreshTokenExpiry int64  `json:"refreshTokenExpiry"`
}

type CreateUserInput struct {
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Roles     []string   `json:"roles"`
	Phone     *string    `json:"phone"`
	CanSms    *bool      `json:"canSMS"`
	Birthday  *time.Time `json:"birthday"`
}

type CreateUserOutput struct {
	ID string `json:"id"`
}

type DeleteUserInput struct {
	ID string `json:"id"`
}

type DeleteUserOutput struct {
	Success bool `json:"success"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserInput struct {
	ID          string     `json:"id"`
	Email       *string    `json:"email"`
	Password    *string    `json:"password"`
	FirstName   *string    `json:"firstName"`
	LastName    *string    `json:"lastName"`
	AddRoles    []string   `json:"addRoles"`
	RemoveRoles []string   `json:"removeRoles"`
	Phone       *string    `json:"phone"`
	CanSms      *bool      `json:"canSMS"`
	Birthday    *time.Time `json:"birthday"`
}

type User struct {
	ID        string     `json:"id"`
	Email     string     `json:"email"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Roles     []string   `json:"roles"`
	Phone     *string    `json:"phone"`
	CanSms    *bool      `json:"canSMS"`
	Birthday  *time.Time `json:"birthday"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	CreatedBy string     `json:"createdBy"`
	UpdatedBy string     `json:"updatedBy"`
}

type UserInput struct {
	ID string `json:"id"`
}
