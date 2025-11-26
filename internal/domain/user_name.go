package domain

import (
	"errors"
)

type UserName string

func NewUserName(name string) (UserName, error) {
	uName := ExistingUserName(name)
	if err := uName.Validate(); err != nil {
		return "", err
	}

	return uName, nil
}

func ExistingUserName(name string) UserName {
	return UserName(name)
}

func (n UserName) Value() string {
	return string(n)
}

func (n UserName) String() string {
	return string(n)
}

func (n UserName) Validate() error {
	if len(n.String()) == 0 {
		return errors.New("user name cannot be empty")
	}

	return nil
}
