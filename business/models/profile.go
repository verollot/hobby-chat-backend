package models

import "github.com/google/uuid"

type User struct {
	ID         string  `json:"user_id"`
	LastName   string  `json:"last_name"`
	FirsttName string  `json:"first_name"`
	Hobbies    []Hobby `json:"hobbies"`
}

type Hobby struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
