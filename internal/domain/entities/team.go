package entities

import (
	"errors"
	"fmt"
	"slices"

	"github.com/alphameo/pr-reviewnager/internal/domain/dto"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
)

const avgUserCountInTeam = 10

type Team struct {
	id   v.ID
	name string
	// slice (not a map) becuse member count cannot be very large
	userIDs []v.ID
}

func NewTeam(name string) (*Team, error) {
	return &Team{
		id:      v.NewID(),
		name:    name,
		userIDs: make([]v.ID, 0, avgUserCountInTeam),
	}, nil
}

func NewExistingTeam(team *dto.TeamDTO) (*Team, error) {
	if team == nil {
		return nil, errors.New("dto cannot be nil")
	}

	err := validateIDsUniqueness(team.UserIDs)
	if err != nil {
		return nil, fmt.Errorf("team members: %w", err)
	}

	userIDs := make([]v.ID, 0, max(len(team.UserIDs), avgUserCountInTeam))
	userIDs = append(userIDs, team.UserIDs...)

	return &Team{
		id:      team.ID,
		name:    team.Name,
		userIDs: userIDs,
	}, nil
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

func (t *Team) AddUser(userID v.ID) error {
	if slices.Contains(t.userIDs, userID) {
		return fmt.Errorf("user with id=%v already exists in team", userID)
	}

	t.userIDs = append(t.userIDs, userID)

	return nil
}

func (t *Team) RemoveUser(userID v.ID) error {
	idx := slices.Index(t.userIDs, userID)
	if idx == -1 {
		return fmt.Errorf("no user with id=%v inside user list", userID)
	}

	t.userIDs = slices.Delete(t.userIDs, idx, idx+1)

	return nil
}
