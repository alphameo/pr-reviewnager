// Package services provides domain services for domain model
package services

import (
	"errors"
	"fmt"
	"math/rand"
	"slices"

	e "github.com/alphameo/pr-reviewnager/internal/domain/entities"
	r "github.com/alphameo/pr-reviewnager/internal/domain/repositories"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
)

type PullRequestDomainService interface {
	CreatePullRequestAndAssignReviewers(pullRequest *e.PullRequest) error
	ReassignReviewerWithUserID(userID v.ID, pullRequestID v.ID) error
	MarkPullRequestAsMergedByIDAndGet(pullRequestID v.ID) (*e.PullRequest, error)
}

type DefaultPullRequestDomainService struct {
	userRepo     r.UserRepository
	teamRepo     r.TeamRepository
	userTeamRepo r.UserTeamRepository
	prRepo       r.PullRequestRepository
}

func NewDefaultPullRequestDomainService(
	userRepository *r.UserRepository,
	pullRequestRepository *r.PullRequestRepository,
	teamRepository *r.TeamRepository,
) (*DefaultPullRequestDomainService, error) {
	if userRepository == nil {
		return nil, errors.New("userRepository cannot be nil")
	}
	if pullRequestRepository == nil {
		return nil, errors.New("pullRequestRepository cannot be nil")
	}
	if teamRepository == nil {
		return nil, errors.New("teamRepository cannot be nil")
	}
	s := DefaultPullRequestDomainService{
		userRepo: *userRepository,
		prRepo:   *pullRequestRepository,
		teamRepo: *teamRepository,
	}
	return &s, nil
}

func (s *DefaultPullRequestDomainService) CreatePullRequestAndAssignReviewers(pullRequest *e.PullRequest) error {
	authorID := pullRequest.AuthorID()
	_, err := s.userRepo.FindByID(authorID)
	if err != nil {
		return err
	}
	team, err := s.userTeamRepo.FindTeamByTeammateID(authorID)
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

func (s *DefaultPullRequestDomainService) ReassignReviewerWithUserID(userID v.ID, pullRequestID v.ID) error {
	pullRequest, err := s.prRepo.FindByID(pullRequestID)
	if err != nil {
		return err
	}
	reviewerIDs := pullRequest.ReviewerIDs()
	reviewerIdx := slices.Index(reviewerIDs, userID)
	if reviewerIdx == -1 {
		return fmt.Errorf("cannot reassign reviewer with id=%s: he is not a reviewer", userID.String())
	}
	authorID := pullRequest.AuthorID()
	team, err := s.userTeamRepo.FindTeamByTeammateID(authorID)
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

	pullRequest.UnassignReviewer(userID)
	pullRequest.AssignReviewer(newReviewer.ID())
	s.prRepo.Update(pullRequest)
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

func (s *DefaultPullRequestDomainService) MarkPullRequestAsMergedByIDAndGet(pullRequestID v.ID) (*e.PullRequest, error) {
	pr, err := s.prRepo.FindByID(pullRequestID)
	if err != nil {
		return nil, err
	}
	if pr == nil {
		return nil, fmt.Errorf("no such pull request with id=%s", pullRequestID.String())
	}

	pr.MarkAsMerged()
	err = s.prRepo.Update(pr)
	if err != nil {
		return nil, err
	}
	return pr, nil
}
