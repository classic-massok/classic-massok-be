package core

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	idField        = "_id"
	contactIDField = "contactID"
	createdAtField = "createdAt"
)

type ID struct {
	OID primitive.ObjectID `bson:"_id"`
	ID  string             `bson:"-"`
}

func (i *ID) SetOID(id primitive.ObjectID) {
	i.OID = id
}

func (i *ID) SetID() {
	i.ID = i.OID.Hex()
}

func (i *ID) GetID() string {
	return i.ID
}

type Accounting struct {
	// Datetime of Entry creation time
	CreatedAt time.Time `bson:"createdAt"`
	// Datetime of last Entry update
	UpdatedAt time.Time `bson:"updatedAt"`

	// TODO: userID of creator/updater
	CreatedBy string `bson:"createdBy,omitempty"`
	UpdatedBy string `bson:"updatedBy,omitempty"`
}

func (a *Accounting) SetAccounting(t time.Time, userID primitive.ObjectID) {
	a.CreatedAt = t
	a.UpdatedAt = t
}
