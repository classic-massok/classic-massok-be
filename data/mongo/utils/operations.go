package utils

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// InsertOne wraps the mongo driver `InsertOne` method
func InsertOne(ctx context.Context, coll *mongo.Collection, doc Creator, loggedInUserID primitive.ObjectID) (*mongo.InsertOneResult, error) {
	doc.SetID(primitive.NewObjectID())
	doc.SetAccounting(time.Now(), loggedInUserID)

	result, err := coll.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}

	return result, nil
}
