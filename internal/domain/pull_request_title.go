package domain

import (
	"errors"
	"strings"
)

type PRTitle string

func NewPRTitle(title string) (PRTitle, error) {
	processed := strings.TrimSpace(title)

	prTitle := ExistingPRTitle(processed)
	if err := prTitle.Validate(); err != nil {
		return "", err
	}

	return prTitle, nil
}

func ExistingPRTitle(title string) PRTitle {
	return PRTitle(title)
}

func (n PRTitle) Value() string {
	return string(n)
}

func (n PRTitle) String() string {
	return string(n)
}

func (n PRTitle) Validate() error {
	if len(n.String()) == 0 {
		return errors.New("pull request title cannot be empty")
	}

	return nil
}
