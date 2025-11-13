// Package services provides domain services for domain model
package services

import (
	"fmt"
	"math/rand"
	"slices"

	e "github.com/alphameo/pr-reviewnager/internal/domain/entities"
	r "github.com/alphameo/pr-reviewnager/internal/domain/repositories"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
)

type PullRequestDomainService interface {
	CreatePullRequest(pullRequest e.PullRequest)
	ReassignReviewer(userID v.ID)
}

type DefaultPullRequestDomainService struct {
	userRepo r.UserRepository
	teamRepo r.TeamRepository
	prRepo   r.PullRequestRepository
}

func NewDefaultPullRequestDomainService(
	userRepository *r.UserRepository,
	pullRequestRepository *r.PullRequestRepository,
	teamRepository *r.TeamRepository,
) *DefaultPullRequestDomainService {
	s := DefaultPullRequestDomainService{
		userRepo: *userRepository,
		prRepo:   *pullRequestRepository,
		teamRepo: *teamRepository,
	}
	return &s
}

func (s *DefaultPullRequestDomainService) CreatePullRequest(pullRequest e.PullRequest) error {
	authorID := pullRequest.AuthorID()
	_, err := s.userRepo.FindById(authorID)
	if err != nil {
		return err
	}
	team, err := s.teamRepo.FindTeamByTeammateID(authorID)
	if err != nil {
		return err
	}
	availableUsers, err := s.userRepo.FindActiveUsersByTeamID(team.ID())
	if err != nil {
		return err
	}
	reviewers := chooseRandomUsers(availableUsers, e.MaxCountOfReviewers, authorID)

	for _, u := range reviewers {
		pullRequest.AssignReviewer(u.ID())
	}
	err = s.prRepo.Create(pullRequest)
	if err != nil {
		return err
	}
	return nil
}

func (s *DefaultPullRequestDomainService) ReassignReviewer(reviewerID v.ID, pullRequestID v.ID) error {
	pullRequest, err := s.prRepo.FindById(pullRequestID)
	if err != nil {
		return err
	}
	reviewerIDs := pullRequest.ReviewerIDs()
	reviewerIdx := slices.Index(reviewerIDs, reviewerID)
	if reviewerIdx == -1 {
		return fmt.Errorf("cannot reassign reviewer with id=%s: he is not a reviewer", reviewerID.String())
	}
	authorID := pullRequest.AuthorID()
	team, err := s.teamRepo.FindTeamByTeammateID(authorID)
	if err != nil {
		return err
	}
	availableUsers, err := s.userRepo.FindActiveUsersByTeamID(team.ID())
	if err != nil {
		return err
	}
	exceptionalReviewerIDs := pullRequest.ReviewerIDs()
	exceptionalReviewerIDs = append(exceptionalReviewerIDs, authorID)
	newReviewer := chooseRandomUser(availableUsers, exceptionalReviewerIDs...)

	pullRequest.UnassignReviewer(reviewerID)
	pullRequest.AssignReviewer(newReviewer.ID())
	s.prRepo.Update(*pullRequest)
	return nil
}

func chooseRandomUser(availableUsers []e.User, except ...v.ID) e.User {
	for i, u := range availableUsers {
		for _, exceptionalU := range except {
			if u.ID() == exceptionalU {
				availableUsers = slices.Delete(availableUsers, i, i+1)
				continue
			}
		}
	}
	idx := rand.Intn(len(availableUsers))
	return availableUsers[idx]
}

func chooseRandomUsers(availableUsers []e.User, maxCount int, except ...v.ID) []e.User {
	for i, u := range availableUsers {
		for _, exceptionalU := range except {
			if u.ID() == exceptionalU {
				availableUsers = slices.Delete(availableUsers, i, i+1)
				continue
			}
		}
	}

	reviewers := make([]e.User, 0, maxCount)
	if len(availableUsers) <= maxCount {
		for i := range min(maxCount, len(availableUsers)) {
			reviewers = append(reviewers, availableUsers[i])
		}
	} else {
		for range maxCount {
			idx := rand.Intn(len(availableUsers))
			reviewers = append(reviewers, availableUsers[idx])
			availableUsers = slices.Delete(availableUsers, idx, idx+1)
		}
	}

	return reviewers
}
