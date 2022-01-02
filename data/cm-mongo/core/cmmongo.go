package core

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Creator interface {
	SetOID(primitive.ObjectID)
	SetAccounting(t time.Time, userID primitive.ObjectID)
}

type Reader interface {
	SetID()
}

type BatchReader interface {
	SetIDs()
}

type Updater interface {
	GetUpdate() bson.M
}

type Order int

const (
	Ascending  Order = 1
	Descending Order = -1
)
