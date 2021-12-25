package cmmongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewCollection(coll *mongo.Collection) *Collection {
	return &Collection{coll}
}

type Collection struct {
	coll *mongo.Collection
}

// Insert wraps the mongo driver `Insert` method
// 	- calls the creator methods to set id and accounting fields
func (c *Collection) Insert(ctx context.Context, loggedInUserID string, doc Creator, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	doc.SetOID(getLoggedInUserOID(loggedInUserID))

	result, err := c.coll.InsertOne(ctx, doc, opts...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Get wraps the mongo driver `FindOne` method
//	- filters by provided id
// 	- fills the provided interface with the returned result
func (c *Collection) Get(ctx context.Context, id string, toFill Reader, opts ...*options.FindOneOptions) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("error converting id string %q to object id: %w", id, err)
	}

	defer toFill.SetID()
	return c.coll.FindOne(ctx, bson.M{idField: oid}, opts...).Decode(toFill)
}

// GetAll wraps the mongo driver `Find` method
//	- filter can be provided
// 	- filsl the provided interface with the returned result
func (c *Collection) GetAll(ctx context.Context, filter bson.M, toFill BatchReader, opts ...*options.FindOptions) error {
	cur, err := c.coll.Find(ctx, filter, opts...)
	if err != nil {
		return err
	}

	defer toFill.SetIDs()
	return cur.All(ctx, toFill)
}

// Edit wraps the mongo driver `FindOneAndUpdate` method
//	- filters by provided id
//	- fills the provided interface with the returned result
func (c *Collection) Edit(ctx context.Context, id, loggedInUserID string, update Updater, toFill Reader, opts ...*options.FindOneAndUpdateOptions) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("error converting id string %q to object id: %w", id, err)
	}

	after := options.After
	opts = append(opts, &options.FindOneAndUpdateOptions{ReturnDocument: &after})

	defer toFill.SetID()
	return c.coll.FindOneAndUpdate(ctx, bson.M{idField: oid}, setUpdates(update.GetUpdate(), loggedInUserID), opts...).Decode(toFill)
}

// Delete wraps the mongo driver `FindOneAndDelete` method
//	- filters by provided id
func (c *Collection) Delete(ctx context.Context, id, loggedInUserID string, opts ...*options.FindOneAndDeleteOptions) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("error converting id string %q to object id: %w", id, err)
	}

	return c.coll.FindOneAndDelete(ctx, bson.M{idField: oid}, opts...).Err()
}

func getLoggedInUserOID(loggedInUserID string) primitive.ObjectID {
	if loggedInUserID == "" {
		return primitive.NilObjectID
	}

	loggedInUserOID, err := primitive.ObjectIDFromHex(loggedInUserID)
	if err != nil {
		return primitive.NilObjectID
		// LOG THIS fmt.Errorf("error converting loggedInUserID string %q to object id: %w", loggedInUserID, err)
	}

	return loggedInUserOID
}

func setUpdates(update bson.M, loggedInUserID string) bson.A {
	update["updatedAt"] = time.Now()
	update["updatedBy"] = getLoggedInUserOID(loggedInUserID)
	return bson.A{bson.M{"$set": update}}
}
