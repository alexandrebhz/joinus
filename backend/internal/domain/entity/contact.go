package entity

import "time"

type Contact struct {
	ID        string
	Name      string
	Email     string
	Subject   string
	Message   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
