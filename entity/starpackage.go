package entity

import "time"

type StarPackage struct {
	ID        uint64
	Amount    Amount
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
