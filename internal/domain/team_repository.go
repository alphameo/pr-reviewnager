package domain

type TeamRepository interface {
	Repository[Team, ID]
	FindByName(teamName string) (*Team, error)
	CreateTeamAndModifyUsers(team *Team, users []*User) error
	FindTeamByTeammateID(userID ID) (*Team, error)
	FindActiveUsersByTeamID(teamID ID) ([]*User, error)
}
