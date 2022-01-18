package models

import (
	"time"

	"github.com/classic-massok/classic-massok-be/lib"
)

const userRole ResourceRole = "users.%s.%s"

type User struct {
	ID        string
	CusKeys   map[string]string
	Email     string
	FirstName string
	LastName  string
	Roles     // TODO: figure out if this has the potential to be nil
	Phone     *string
	CanSMS    *bool
	Birthday  *time.Time
	Accounting
}

func (u *User) acl() ACL {
	return ACL{
		{
			Roles: Roles{
				userRole.Populate(RoleTypeUser, u.ID),
			},
			Actions: lib.NewStringset(
				"user.read",
				"user.update",
			),
		},
	}
}

type UserEdit struct { // TODO: need to figure out adding/removing roles
	Email     *string
	Password  *string
	FirstName *string
	LastName  *string
	Roles
	Phone    *string
	CanSMS   *bool
	Birthday *time.Time
}
