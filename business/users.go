package business

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/classic-massok/classic-massok-be/data/mongo"
	"github.com/classic-massok/classic-massok-be/data/mongo/cmmongo"
	"github.com/classic-massok/classic-massok-be/lib"
	"github.com/labstack/gommon/random"
	"golang.org/x/crypto/bcrypt"
)

// NewUsersBiz is the constructure for the users business layers
func NewUsersBiz() *usersBiz {
	data, err := mongo.NewUsersData(cmmongo.Database)
	if err != nil {
		panic(fmt.Sprintf("error constructing users biz: %v", err))
	}

	return &usersBiz{data}
}

// usersBiz is the business layer for interacting with user data
type usersBiz struct {
	data usersData
}

func (u *usersBiz) Authn(ctx context.Context, email, password string) (string, string, error) {
	mongoUser, err := u.data.GetByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}

	return mongoUser.GetID(), mongoUser.CusKey, bcrypt.CompareHashAndPassword(mongoUser.Password, []byte(password))
}

func (u *usersBiz) New(ctx context.Context, loggedInUserID string, password string, user User) (string, error) {
	if !user.Roles.Validate() {
		return "", fmt.Errorf("invalid user roles")
	}

	passwordBytes := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return u.data.New(ctx, loggedInUserID, mongo.User{
		CusKey:    random.String(15),
		Email:     user.Email,
		Password:  hashedPassword,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		CanSMS:    user.CanSMS,
		Birthday:  user.Birthday,
	})
}

func (u *usersBiz) Get(ctx context.Context, id string) (*User, error) {
	mongoUser, err := u.data.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &User{
		mongoUser.GetID(),
		mongoUser.CusKey,
		mongoUser.Email,
		mongoUser.FirstName,
		mongoUser.LastName,
		mongoUser.Roles,
		mongoUser.Phone,
		mongoUser.CanSMS,
		mongoUser.Birthday,
		accounting{
			mongoUser.CreatedAt,
			mongoUser.UpdatedAt,
			mongoUser.CreatedBy,
			mongoUser.UpdatedBy,
		},
	}, nil
}

func (u *usersBiz) GetAll(ctx context.Context) ([]*User, error) {
	mongoUsers, err := u.data.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]*User, len(mongoUsers))
	for i, mongoUser := range mongoUsers {
		users[i] = &User{
			mongoUser.GetID(),
			mongoUser.CusKey,
			mongoUser.Email,
			mongoUser.FirstName,
			mongoUser.LastName,
			mongoUser.Roles,
			mongoUser.Phone,
			mongoUser.CanSMS,
			mongoUser.Birthday,
			accounting{
				mongoUser.CreatedAt,
				mongoUser.UpdatedAt,
				mongoUser.CreatedBy,
				mongoUser.UpdatedBy,
			},
		}
	}

	return users, nil
}

func (u *usersBiz) Edit(ctx context.Context, id, loggedInUserID string, updateCusKey bool, userEdit UserEdit) (*User, error) {
	mongoUserEdit := mongo.UserEdit{
		Email:     userEdit.Email,
		FirstName: userEdit.FirstName,
		LastName:  userEdit.LastName,
		Phone:     userEdit.Phone,
		CanSMS:    userEdit.CanSMS,
		Birthday:  userEdit.Birthday,
	}

	passwordBytes := []byte{}
	if userEdit.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		mongoUserEdit.Password = &hashedPassword
	} else {
		mongoUserEdit.Password = nil
	}

	if userEdit.Password != nil || updateCusKey {
		newCusKey := random.String(15)
		mongoUserEdit.CusKey = &newCusKey
	}

	if userEdit.Roles != nil {
		if !userEdit.Roles.Validate() {
			return nil, fmt.Errorf("invalid user roles")
		}

		mongoUserEdit.Roles = (*[]string)(&userEdit.Roles)
	}

	mongoUser, err := u.data.Edit(ctx, id, loggedInUserID, mongoUserEdit)
	if err != nil {
		return nil, err
	}

	return &User{
		mongoUser.GetID(),
		mongoUser.CusKey,
		mongoUser.Email,
		mongoUser.FirstName,
		mongoUser.LastName,
		mongoUser.Roles,
		mongoUser.Phone,
		mongoUser.CanSMS,
		mongoUser.Birthday,
		accounting{
			mongoUser.CreatedAt,
			mongoUser.UpdatedAt,
			mongoUser.CreatedBy,
			mongoUser.UpdatedBy,
		},
	}, nil
}

func (u *usersBiz) Delete(ctx context.Context, id, loggedInUserID string) error {
	return u.data.Delete(ctx, id, loggedInUserID)
}

type User struct {
	id        string
	cusKey    string
	Email     string
	FirstName string
	LastName  string
	Roles     // TODO: figure out if this has the potential to be nil
	Phone     *string
	CanSMS    *bool
	Birthday  *time.Time
	accounting
}

func (u *User) acl() ACL {
	return ACL{
		{
			Roles: Roles{
				roleAdmin.String(),
			},
			Actions: lib.NewStringset(
				"",
			),
		},
	}
}

func (u *User) GetID() string {
	return u.id
}

func (u *User) GetCusKey() string {
	return u.cusKey
}

type UserEdit struct {
	Email     *string
	Password  *string
	FirstName *string
	LastName  *string
	Roles
	Phone    *string
	CanSMS   *bool
	Birthday *time.Time
}

type Roles []string

func (r Roles) SetRole(appScope ApplicationScope, role Role) {
	if r == nil {
		r = Roles{}
	}

	prefixedRole := fmt.Sprintf("%s.%s", appScope, role)
	r = append(r, prefixedRole)
}

func (r Roles) HasRole(appScope ApplicationScope, role Role) bool {
	if r == nil {
		return false
	}

	prefixedRole := fmt.Sprintf("%s.%s", appScope, role)
	for _, curRole := range r {
		if curRole == prefixedRole {
			return true
		}
	}

	return false
}

func (r Roles) RemoveRole(appScope ApplicationScope, role Role) bool {
	if r == nil {
		return false
	}

	prefixedRole := fmt.Sprintf("%s.%s", appScope, role)
	for i, curRole := range r {
		if curRole == prefixedRole {
			r = append(r[:i], r[i+1:]...)
			return true
		}
	}

	return false
}

func (r Roles) Validate() bool {
	if r == nil {
		return true
	}

	for _, curRole := range r {
		rs := strings.Split(curRole, ".")
		if len(rs) != 2 {
			return false
		}

		if !ApplicationScope(rs[0]).Validate() {
			return false
		}

		if !Role(rs[1]).Validate() {
			return false
		}
	}

	return true
}

//counterfeiter:generate . usersData
type usersData interface {
	New(ctx context.Context, loggedInUserID string, user mongo.User) (string, error)
	Get(ctx context.Context, id string) (*mongo.User, error)
	GetByEmail(ctx context.Context, email string) (*mongo.User, error)
	GetAll(ctx context.Context) ([]*mongo.User, error)
	Edit(ctx context.Context, id, loggedInUserID string, edit mongo.UserEdit) (*mongo.User, error)
	Delete(ctx context.Context, id, loggedInUserID string) error
}

type userGetter interface {
	Get(ctx context.Context, id string) (*User, error)
}
