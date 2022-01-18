package models

import "time"

type Accounting struct {
	// Unix timestamp of Entry creation/update time
	CreatedAt time.Time
	UpdatedAt time.Time

	// UserID of creator/updater
	CreatedBy string
	UpdatedBy string
}
