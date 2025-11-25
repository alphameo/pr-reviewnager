package domain

type TeamRepository interface {
	Repository[Team, TeamDTO, ID]
	FindByName(teamName string) (*TeamDTO, error)
	CreateTeamAndModifyUsers(team *Team, users []*User) error
	FindTeamByTeammateID(userID ID) (*TeamDTO, error)
	FindActiveUsersByTeamID(teamID ID) ([]*UserDTO, error)
}
