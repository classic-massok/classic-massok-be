package business

import (
	"context"
	"fmt"
	"time"

	cmmongo "github.com/classic-massok/classic-massok-be/data/cm-mongo"
	"github.com/classic-massok/classic-massok-be/lib"
	"github.com/labstack/gommon/random"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

const userResource ResourceRole = "users.%s.%s"

const (
	cusKeysKey  = "CusKeys"
	ipAdressKey = "IPAddress"
)

// NewUsersBiz is the constructure for the users business layers
func NewUsersBiz(db *mongo.Database) *usersBiz {
	data, err := cmmongo.NewUsersData(db)
	if err != nil {
		panic(fmt.Sprintf("error constructing users biz: %v", err))
	}

	return &usersBiz{data}
}

// usersBiz is the business layer for interacting with user data
type usersBiz struct {
	data usersData
}

func (u *usersBiz) Authn(ctx context.Context, email, password string) (string, map[string]string, error) {
	mongoUser, err := u.data.GetByEmail(ctx, email)
	if err != nil {
		return "", nil, err
	}

	return mongoUser.GetID(), mongoUser.CusKeys, bcrypt.CompareHashAndPassword(mongoUser.Password, []byte(password))
}

func (u *usersBiz) New(ctx context.Context, loggedInUserID string, password string, user User) (string, error) {
	if user.Roles == nil {
		user.Roles = Roles{}
	}

	user.Roles = append(user.Roles, "users.user.self")
	user.Roles.DeDupe()

	if err := user.Roles.Validate(); err != nil {
		return "", err
	}

	passwordBytes := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return u.data.New(ctx, loggedInUserID, cmmongo.User{
		CusKeys:   map[string]string{ctx.Value(ipAdressKey).(string): random.String(15)},
		Email:     user.Email,
		Password:  hashedPassword,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Roles:     user.Roles,
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
		mongoUser.CusKeys,
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
			mongoUser.CusKeys,
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
	mongoUserEdit := cmmongo.UserEdit{
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

		mongoUserEdit.Password = hashedPassword
	} else {
		mongoUserEdit.Password = nil
	}

	if userEdit.Password != nil || updateCusKey {
		mongoUserEdit.CusKeys = ctx.Value(cusKeysKey).(map[string]string)
		mongoUserEdit.CusKeys[ctx.Value(ipAdressKey).(string)] = random.String(15)
	}

	if userEdit.Roles != nil {
		if err := userEdit.Roles.Validate(); err != nil {
			return nil, err
		}

		userEdit.Roles.DeDupe()
		mongoUserEdit.Roles = ([]string)(userEdit.Roles)
	} else {
		userEdit.Roles = []string{}
	}

	mongoUser, err := u.data.Edit(ctx, id, loggedInUserID, mongoUserEdit)
	if err != nil {
		return nil, err
	}

	return &User{
		mongoUser.GetID(),
		mongoUser.CusKeys,
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
	cusKeys   map[string]string
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
				userResource.Populate(roleTypeUser, u.id),
			},
			Actions: lib.NewStringset(
				"user.read",
				"user.update",
			),
		},
	}
}

func (u *User) GetID() string {
	return u.id
}

func (u *User) GetCusKey(ipAdress string) string {
	return u.cusKeys[ipAdress]
}

func (u *User) GetCusKeys() map[string]string {
	return u.cusKeys
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

//counterfeiter:generate . usersData
type usersData interface {
	New(ctx context.Context, loggedInUserID string, user cmmongo.User) (string, error)
	Get(ctx context.Context, id string) (*cmmongo.User, error)
	GetByEmail(ctx context.Context, email string) (*cmmongo.User, error)
	GetAll(ctx context.Context) ([]*cmmongo.User, error)
	Edit(ctx context.Context, id, loggedInUserID string, edit cmmongo.UserEdit) (*cmmongo.User, error)
	Delete(ctx context.Context, id, loggedInUserID string) error
}

type userGetter interface {
	Get(ctx context.Context, id string) (*User, error)
}
