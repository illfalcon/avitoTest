package models

import "time"

type User struct {
	ID        int
	Username  string
	CreatedAt time.Time
}

type Chat struct {
	ID        int
	Name      string
	Users     []User
	CreatedAt time.Time
}

type Message struct {
	ID        int
	Chat      int
	Author    int
	Text      string
	CreatedAt time.Time
}
