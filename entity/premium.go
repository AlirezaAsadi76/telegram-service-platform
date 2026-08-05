package entity

import "time"

type PremiumPlan struct {
	ID        uint64
	Months    uint8
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
