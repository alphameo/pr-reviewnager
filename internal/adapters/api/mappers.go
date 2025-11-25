package api

import (
	"time"

	"github.com/alphameo/pr-reviewnager/internal/app"
	"github.com/alphameo/pr-reviewnager/internal/domain"
)

func ToAPITeam(d app.TeamWithUsersDTO) Team {
	members := make([]TeamMember, len(d.TeamUsers))
	for i, m := range d.TeamUsers {
		member := ToAPITeamMember(*m)
		members[i] = member
	}

	return Team{
		TeamName: d.TeamName,
		Members:  members,
	}
}

func FromAPITeam(t Team) app.TeamWithUsersDTO {
	members := make([]*domain.UserDTO, len(t.Members))
	for i, m := range t.Members {
		member := FromAPITeamMember(m)
		members[i] = &member
	}

	return app.TeamWithUsersDTO{
		TeamName:  t.TeamName,
		TeamUsers: members,
	}
}

func ToAPITeamMember(m domain.UserDTO) TeamMember {
	return TeamMember{
		UserId:   m.ID.String(),
		Username: m.Name,
		IsActive: m.Active,
	}
}

func FromAPITeamMember(m TeamMember) domain.UserDTO {
	id, _ := domain.NewIDFromString(m.UserId)

	return domain.UserDTO{
		ID:     id,
		Name:   m.Username,
		Active: m.IsActive,
	}
}

func ToAPIUser(u app.UserWithTeamNameDTO) User {
	return User{
		UserId:   u.User.ID.String(),
		Username: u.User.Name,
		TeamName: u.TeamName,
		IsActive: u.User.Active,
	}
}

func ToAPIPullRequest(d domain.PullRequestDTO) PullRequest {
	reviewers := make([]string, len(d.ReviewerIDs))
	for i, rid := range d.ReviewerIDs {
		reviewers[i] = rid.String()
	}

	var mergedAt *time.Time
	if d.MergedAt != nil {
		mergedAt = d.MergedAt
	}

	return PullRequest{
		PullRequestId:     d.ID.String(),
		PullRequestName:   d.Title,
		AuthorId:          d.AuthorID.String(),
		Status:            PullRequestStatus(d.Status), // "OPEN" / "MERGED"
		AssignedReviewers: reviewers,
		CreatedAt:         &d.CreatedAt,
		MergedAt:          mergedAt,
	}
}

func ToAPIPullRequestShort(d domain.PullRequestDTO) PullRequestShort {
	return PullRequestShort{
		PullRequestId:   d.ID.String(),
		PullRequestName: d.Title,
		AuthorId:        d.AuthorID.String(),
		Status:          PullRequestShortStatus(d.Status),
	}
}

func ToAPIPullRequestShortList(list []*domain.PullRequestDTO) []PullRequestShort {
	out := make([]PullRequestShort, len(list))
	for i, pr := range list {
		out[i] = ToAPIPullRequestShort(*pr)
	}
	return out
}
