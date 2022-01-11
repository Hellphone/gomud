package models

import "time"

const (
	StatusOffline Status = iota
	StatusOnline
)

type Status int

type User struct {
	Login          string
	Level          int
	Status         Status
	LastActionTime time.Time
	Location       string
}

// TODO: change status in DB (no need)
func (u *User) SwitchStatus() {
	if u.Status == StatusOffline {
		u.Status = StatusOnline
	} else {
		u.Status = StatusOffline
	}
}
