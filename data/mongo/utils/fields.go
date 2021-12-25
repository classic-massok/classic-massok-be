package utils

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ID struct {
	ID primitive.ObjectID `bson:"_id"`
}

func (i *ID) SetID(id primitive.ObjectID) {
	i.ID = id
}

func (i *ID) GetID() primitive.ObjectID { return i.ID }

type Accounting struct {
	// Unix timestamp of Entry creation time
	CreatedAt time.Time `bson:"createdAt"`
	// Unix timestamp of last Entry update
	UpdatedAt time.Time `bson:"updatedAt"`

	// TODO: userID of creator/updater
	CreatedBy string `bson:"createdBy,omitempty"`
	UpdatedBy string `bson:"updatedBy,omitempty"`
}

func (a *Accounting) SetAccounting(t time.Time, userID primitive.ObjectID) {
	a.CreatedAt = t
	a.UpdatedAt = t
}

func (a *Accounting) RecordUpdate(t time.Time) {
	a.UpdatedAt = t
}

func (a *Accounting) GetCreatedAt() time.Time { return a.CreatedAt }
func (a *Accounting) GetUpdatedAt() time.Time { return a.UpdatedAt }
