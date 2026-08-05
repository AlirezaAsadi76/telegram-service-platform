package entity

import "time"

type AdsPlan struct {
	ID        uint64
	Views     uint64
	CPM       float64
	Price     float64
	DailyShow uint8
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
