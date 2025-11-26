package domain

import (
	"errors"
	"fmt"
	"slices"
)

const avgUserCountInTeam = 10

var ErrAlreadyTeamMember = errors.New("user is already a team member")

type Team struct {
	id   ID
	name TeamName
	// slice (not a map) becuse member count cannot be very large
	userIDs []ID
}

func NewTeam(name TeamName) (*Team, error) {
	return &Team{
		id:      NewID(),
		name:    name,
		userIDs: make([]ID, 0, avgUserCountInTeam),
	}, nil
}

func ExistingTeam(
	id ID,
	name TeamName,
	userIDs []ID,
) *Team {
	uIDs := make([]ID, 0, max(len(userIDs), avgUserCountInTeam))
	uIDs = append(uIDs, userIDs...)

	return &Team{
		id:      id,
		name:    name,
		userIDs: uIDs,
	}
}

func (t *Team) ID() ID {
	return t.id
}

func (t *Team) Name() TeamName {
	return t.name
}

func (t *Team) UserIDs() []ID {
	return slices.Clone(t.userIDs)
}

func (t *Team) AddUser(userID ID) error {
	if slices.Contains(t.userIDs, userID) {
		return fmt.Errorf("%w: userID=%v", ErrAlreadyTeamMember, userID)
	}

	t.userIDs = append(t.userIDs, userID)

	return nil
}

func (t *Team) RemoveUser(userID ID) error {
	idx := slices.Index(t.userIDs, userID)
	if idx == -1 {
		return fmt.Errorf("no user with id=%v inside user list", userID)
	}

	t.userIDs = slices.Delete(t.userIDs, idx, idx+1)

	return nil
}

func (t *Team) Validate() error {
	if err := validateIDsUniqueness(t.UserIDs()); err != nil {
		return fmt.Errorf("error in team members: %w", err)
	}
	if err := t.name.Validate(); err != nil {
		return err
	}

	return nil
}
