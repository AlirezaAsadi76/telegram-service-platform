package entity

import "time"

type StarPackage struct {
	ID        uint64
	Amount    uint64
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
