package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/classic-massok-be/data/mongo/utils"
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
		coll: db.Collection(UsersCollection),
	}, nil
}

// usersData is the data layer that access the users db collection
type usersData struct {
	coll *mongo.Collection
}

// User represents the user db model
type User struct {
	utils.ID `bson:",inline"`

	// Required fields
	Email     string `bson:"email"`
	Password  []byte `bson:"password"` // hashed
	FirstName string `bson:"firstName"`
	LastName  string `bson:"lastName"`

	// Optional fields
	Phone    *string    `bson:"phone"`
	CanSMS   *bool      `bson:"canSMS"`
	Birthday *time.Time `bson:"birthday"`

	utils.Accounting `bson:",inline"`
}

func (u *usersData) New(ctx context.Context, loggedInUserID primitive.ObjectID, user User) (string, error) {
	result, err := utils.InsertOne(ctx, u.coll, &user, loggedInUserID)
	if err != nil {
		return "", fmt.Errorf("error creating new user: %w", err)
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}
