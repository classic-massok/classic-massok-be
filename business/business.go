package business

import "time"

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

type Accounting struct {
	// Unix timestamp of Entry creation time
	createdAt time.Time
	// Unix timestamp of last Entry update
	updatedAt time.Time

	// TODO: userID of creator/updater
	createdBy string
	updatedBy string
}

func (a *Accounting) GetCreatedAt() time.Time {
	return a.createdAt
}

func (a *Accounting) GetUpdatedAt() time.Time {
	return a.updatedAt
}

func (a *Accounting) GetCreatedBy() string {
	return a.createdBy
}

func (a *Accounting) GetUpdatedBy() string {
	return a.updatedBy
}
