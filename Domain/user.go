package domain

import "time"

type User struct {
	UserId    int64
	Role      string
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
