package utils

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var Database = &mongo.Database{}

type Creator interface {
	SetID(primitive.ObjectID)
	SetAccounting(t time.Time, userID primitive.ObjectID)
}

type Updater interface {
	UpdateAccounting(t time.Time, userID primitive.ObjectID)
}

type Order int

const (
	Ascending  Order = 1
	Descending Order = -1
)
