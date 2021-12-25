// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

import (
	"time"
)

type CreateUserInput struct {
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
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

type UpdateUserInput struct {
	ID        string     `json:"id"`
	Email     *string    `json:"email"`
	Password  *string    `json:"password"`
	FirstName *string    `json:"firstName"`
	LastName  *string    `json:"lastName"`
	Phone     *string    `json:"phone"`
	CanSms    *bool      `json:"canSMS"`
	Birthday  *time.Time `json:"birthday"`
}

type User struct {
	ID        string     `json:"id"`
	Email     string     `json:"email"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Phone     *string    `json:"phone"`
	CanSms    *bool      `json:"canSMS"`
	Birthday  *time.Time `json:"birthday"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	CreatedBy string     `json:"createdBy"`
	UpdatedBy string     `json:"updatedBy"`
}
