// Package entities provides domain model of service
package entities

import (
	"errors"

	"github.com/alphameo/pr-reviewnager/internal/domain/dto"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
)

type User struct {
	id     v.ID
	name   string
	active bool
}

func NewExistingUser(user *dto.UserDTO) (*User, error) {
	if user == nil {
		return nil, errors.New("dto cannot be nil")
	}

	return &User{
		id:     user.ID,
		name:   user.Name,
		active: user.Active,
	}, nil
}

func NewUser(name string, active bool) (*User, error) {
	return &User{
		id:     v.NewID(),
		name:   name,
		active: active,
	}, nil
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
