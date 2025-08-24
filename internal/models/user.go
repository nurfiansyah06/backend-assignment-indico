package models

import "time"

type User struct {
	Id        int
	Name      string
	CreatedAt time.Time
}

type RequestUser struct {
	Name string
}
