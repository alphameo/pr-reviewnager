package api

import (
	"time"

	s "github.com/alphameo/pr-reviewnager/internal/application/services"
	"github.com/alphameo/pr-reviewnager/internal/domain/dto"
	"github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
)

func ToAPITeam(d s.TeamWithUsersDTO) Team {
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

func FromAPITeam(t Team) s.TeamWithUsersDTO {
	members := make([]*dto.UserDTO, len(t.Members))
	for i, m := range t.Members {
		member := FromAPITeamMember(m)
		members[i] = &member
	}

	return s.TeamWithUsersDTO{
		TeamName:  t.TeamName,
		TeamUsers: members,
	}
}

func ToAPITeamMember(m dto.UserDTO) TeamMember {
	return TeamMember{
		UserId:   m.ID.String(),
		Username: m.Name,
		IsActive: m.Active,
	}
}

func FromAPITeamMember(m TeamMember) dto.UserDTO {
	id, _ := valueobjects.NewIDFromString(m.UserId)

	return dto.UserDTO{
		ID:     id,
		Name:   m.Username,
		Active: m.IsActive,
	}
}

func ToAPIUser(u s.UserWithTeamNameDTO) User {
	return User{
		UserId:   u.User.ID.String(),
		Username: u.User.Name,
		TeamName: u.TeamName,
		IsActive: u.User.Active,
	}
}

func ToAPIPullRequest(d dto.PullRequestDTO) PullRequest {
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

func ToAPIPullRequestShort(d dto.PullRequestDTO) PullRequestShort {
	return PullRequestShort{
		PullRequestId:   d.ID.String(),
		PullRequestName: d.Title,
		AuthorId:        d.AuthorID.String(),
		Status:          PullRequestShortStatus(d.Status),
	}
}

func ToAPIPullRequestShortList(list []*dto.PullRequestDTO) []PullRequestShort {
	out := make([]PullRequestShort, len(list))
	for i, pr := range list {
		out[i] = ToAPIPullRequestShort(*pr)
	}
	return out
}
