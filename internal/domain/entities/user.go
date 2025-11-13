// Package entities provides domain model of service
package entities

import (
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
)

type User struct {
	id     v.ID
	name   string
	active bool
}

func (u *User) ID() v.ID {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Active() bool {
	return u.active
}

func (u *User) SetActive(active bool) {
	u.active = active
}
