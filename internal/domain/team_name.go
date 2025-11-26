package domain

import (
	"errors"
	"strings"
)

type TeamName string

func NewTeamName(name string) (TeamName, error) {
	processed := strings.TrimSpace(name)

	tName := ExistingTeamName(processed)
	if err := tName.Validate(); err != nil {
		return "", err
	}

	return tName, nil
}

func ExistingTeamName(name string) TeamName {
	return TeamName(name)
}

func (n TeamName) Value() string {
	return string(n)
}

func (n TeamName) String() string {
	return string(n)
}

func (n TeamName) Validate() error {
	if len(n.String()) == 0 {
		return errors.New("team name cannot be empty")
	}

	return nil
}
