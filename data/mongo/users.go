package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/classic-massok/classic-massok-be/data/mongo/cmmongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UsersCollection represents the name for the users collection
const UsersCollection = "users"

// NewUsersData is the constructor for the users data layer
func NewUsersData(db *mongo.Database) (*usersData, error) {
	if db == nil {
		return nil, fmt.Errorf("database cannot be nil")
	}

	return &usersData{
		coll: cmmongo.NewCollection(db.Collection(UsersCollection)),
	}, nil
}

// usersData is the data layer that access the users db collection
type usersData struct {
	coll *cmmongo.Collection
}

// Users represents an array of users
type Users []*User

func (us Users) SetIDs() {
	for _, user := range us {
		user.SetID()
	}
}

// User represents the user db model
type User struct {
	cmmongo.ID `bson:",inline"`

	// Required fields
	Email     string `bson:"email"`
	Password  []byte `bson:"password"` // hashed
	FirstName string `bson:"firstName"`
	LastName  string `bson:"lastName"`

	// Optional fields
	Phone    *string    `bson:"phone"`
	CanSMS   *bool      `bson:"canSMS"`
	Birthday *time.Time `bson:"birthday"`

	cmmongo.Accounting `bson:",inline"`
}

type UserEdit struct {
	Email     *string
	Password  *[]byte
	FirstName *string
	LastName  *string
	Phone     *string
	CanSMS    *bool
	Birthday  *time.Time
}

func (ue *UserEdit) GetUpdate() bson.M {
	update := bson.M{}
	if ue.Email != nil {
		update["email"] = *ue.Email
	}

	if ue.FirstName != nil {
		update["firstName"] = *ue.FirstName
	}

	if ue.LastName != nil {
		update["lastName"] = *ue.LastName
	}

	if ue.Phone != nil {
		update["phone"] = ue.Phone
	}

	if ue.CanSMS != nil {
		update["canSMS"] = ue.CanSMS
	}

	if ue.Birthday != nil {
		update["birthday"] = ue.Birthday
	}

	return update
}

func (u *usersData) New(ctx context.Context, loggedInUserID string, user User) (string, error) {
	result, err := u.coll.Insert(ctx, loggedInUserID, &user)
	if err != nil {
		return "", fmt.Errorf("error creating new user: %w", err)
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (u *usersData) Get(ctx context.Context, id string) (*User, error) {
	var user User
	if err := u.coll.Get(ctx, id, &user); err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return &user, nil
}

func (u *usersData) GetAll(ctx context.Context) ([]*User, error) {
	var users Users
	if err := u.coll.GetAll(ctx, bson.M{}, &users); err != nil {
		return nil, fmt.Errorf("error getting all users: %w", err)
	}

	return users, nil
}

func (u *usersData) Edit(ctx context.Context, id, loggedInuserID string, user UserEdit) (*User, error) {
	var updatedUser User
	if err := u.coll.Edit(ctx, id, loggedInuserID, &user, &updatedUser); err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return &updatedUser, nil
}

func (u *usersData) Delete(ctx context.Context, id, loggedInUserID string) error {
	return u.coll.Delete(ctx, id, loggedInUserID)
}
