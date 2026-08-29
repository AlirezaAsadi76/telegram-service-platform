package entity

import "time"

type StarPackage struct {
	ID        uint64
	Amount    int64
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
