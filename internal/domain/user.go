package domain

import "errors"

type UserDTO struct {
	ID     ID
	Name   string
	Active bool
}
type User struct {
	id     ID
	name   string
	active bool
}

func NewExistingUser(user *UserDTO) (*User, error) {
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
		id:     NewID(),
		name:   name,
		active: active,
	}, nil
}

func (u *User) ID() ID {
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
