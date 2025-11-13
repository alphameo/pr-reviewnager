package entities

import (
	"fmt"
	"slices"

	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
)

const avgUserCountInTeam = 10

type Team struct {
	id      v.ID
	name    string
	userIDs []v.ID
}

func NewTeam(name string) *Team {
	t := Team{
		id:      v.NewID(),
		name:    name,
		userIDs: make([]v.ID, 0, avgUserCountInTeam),
	}

	return &t
}

func (t *Team) ID() v.ID {
	return t.id
}

func (t *Team) Name() string {
	return t.name
}

func (t *Team) UserIDs() []v.ID {
	return slices.Clone(t.userIDs)
}

func (t *Team) AddUser(userID v.ID) {
	t.userIDs = append(t.userIDs, userID)
}

func (t *Team) RemoveUser(userID v.ID) error {
	idx := slices.Index(t.userIDs, userID)
	if idx == -1 {
		return fmt.Errorf("cannot remove user: no user with id=%d inside user list", userID)
	}

	t.userIDs = slices.Delete(t.userIDs, idx, idx+1)
	return nil
}
