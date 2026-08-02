package entity

import "time"

type User struct {
	ID          uint64
	TelegramID  int64
	Username    string
	FirstName   string
	LastName    string
	PhoneNumber string
	Role        Role
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
