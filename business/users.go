package business

import (
	"context"
	"fmt"

	"github.com/classic-massok/classic-massok-be/business/models"
	cmmongo "github.com/classic-massok/classic-massok-be/data/cm-mongo"
	"github.com/classic-massok/classic-massok-be/lib"
	"github.com/labstack/gommon/random"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
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

func (u *usersBiz) New(ctx context.Context, loggedInUserID string, password string, user models.User) (string, error) {
	if user.Roles == nil {
		user.Roles = models.Roles{}
	}

	user.Roles = append(user.Roles, models.UserSelf)
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
		CusKeys:   map[string]string{ctx.Value(lib.IPAddressKey).(string): random.String(15)},
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

func (u *usersBiz) Get(ctx context.Context, id string) (*models.User, error) {
	mongoUser, err := u.data.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.User{
		mongoUser.GetID(),
		mongoUser.CusKeys,
		mongoUser.Email,
		mongoUser.FirstName,
		mongoUser.LastName,
		mongoUser.Roles,
		mongoUser.Phone,
		mongoUser.CanSMS,
		mongoUser.Birthday,
		models.Accounting{
			mongoUser.CreatedAt,
			mongoUser.UpdatedAt,
			mongoUser.CreatedBy,
			mongoUser.UpdatedBy,
		},
	}, nil
}

func (u *usersBiz) GetAll(ctx context.Context) ([]*models.User, error) {
	mongoUsers, err := u.data.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]*models.User, len(mongoUsers))
	for i, mongoUser := range mongoUsers {
		users[i] = &models.User{
			mongoUser.GetID(),
			mongoUser.CusKeys,
			mongoUser.Email,
			mongoUser.FirstName,
			mongoUser.LastName,
			mongoUser.Roles,
			mongoUser.Phone,
			mongoUser.CanSMS,
			mongoUser.Birthday,
			models.Accounting{
				mongoUser.CreatedAt,
				mongoUser.UpdatedAt,
				mongoUser.CreatedBy,
				mongoUser.UpdatedBy,
			},
		}
	}

	return users, nil
}

func (u *usersBiz) Edit(ctx context.Context, id, loggedInUserID string, updateCusKey bool, userEdit models.UserEdit) (*models.User, error) {
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
		mongoUserEdit.CusKeys = ctx.Value(lib.CusKeysKey).(map[string]string)
		mongoUserEdit.CusKeys[ctx.Value(lib.IPAddressKey).(string)] = random.String(15)
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

	return &models.User{
		mongoUser.GetID(),
		mongoUser.CusKeys,
		mongoUser.Email,
		mongoUser.FirstName,
		mongoUser.LastName,
		mongoUser.Roles,
		mongoUser.Phone,
		mongoUser.CanSMS,
		mongoUser.Birthday,
		models.Accounting{
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
	Get(ctx context.Context, id string) (*models.User, error)
}
