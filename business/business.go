package business

import (
	"time"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

type accounting struct {
	// Unix timestamp of Entry creation time
	createdAt time.Time
	// Unix timestamp of last Entry update
	updatedAt time.Time

	// TODO: userID of creator/updater
	createdBy string
	updatedBy string
}

func (a *accounting) GetCreatedAt() time.Time {
	return a.createdAt
}

func (a *accounting) GetUpdatedAt() time.Time {
	return a.updatedAt
}

func (a *accounting) GetCreatedBy() string {
	return a.createdBy
}

func (a *accounting) GetUpdatedBy() string {
	return a.updatedBy
}
