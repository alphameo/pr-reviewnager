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
	// CreateWithReviewers() creates a new pull request and automatically assigns
	// 2 reviewers by randomly selecting from the author's team.
	CreateWithReviewers(pullRequest *e.PullRequest) error

	// ReassignReviewer() unassign user-reviewer with given id and assigns another from his team, excluding
	// him and pr author. After, method returns id of new user-reviewer and pull request
	ReassignReviewer(userID v.ID, pullRequestID v.ID) (*ReassignReviewerResponse, error)

	// MarkAsMerged() idempotently marks pull request as merged and sets time of marking
	MarkAsMerged(pullRequestID v.ID) (*e.PullRequest, error)
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

func (s *DefaultPullRequestDomainService) CreateWithReviewers(pullRequest *e.PullRequest) error {
	authorID := pullRequest.AuthorID()
	_, err := s.userRepo.FindByID(authorID)
	if err != nil {
		return err
	}
	team, err := s.userTeamRepo.FindTeamByTeammateID(authorID)
	if err != nil {
		return err
	}
	availableUsers, err := s.userTeamRepo.FindActiveUsersByTeamID(team.ID())
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

type ReassignReviewerResponse struct {
	NewReviewerID v.ID
	PullRequest   e.PullRequest
}

func (s *DefaultPullRequestDomainService) ReassignReviewer(userID v.ID, pullRequestID v.ID) (*ReassignReviewerResponse, error) {
	pullRequest, err := s.prRepo.FindByID(pullRequestID)
	if err != nil {
		return nil, err
	}
	reviewerIDs := pullRequest.ReviewerIDs()
	reviewerIdx := slices.Index(reviewerIDs, userID)
	if reviewerIdx == -1 {
		return nil, fmt.Errorf("cannot reassign reviewer with id=%s: he is not a reviewer", userID.String())
	}
	authorID := pullRequest.AuthorID()
	team, err := s.userTeamRepo.FindTeamByTeammateID(authorID)
	if err != nil {
		return nil, err
	}
	availableUsers, err := s.userTeamRepo.FindActiveUsersByTeamID(team.ID())
	if err != nil {
		return nil, err
	}
	exceptionalReviewerIDs := pullRequest.ReviewerIDs()
	exceptionalReviewerIDs = append(exceptionalReviewerIDs, authorID)
	newReviewer := chooseRandomUser(availableUsers, exceptionalReviewerIDs...)

	pullRequest.UnassignReviewer(userID)
	pullRequest.AssignReviewer(newReviewer.ID())
	err = s.prRepo.Update(pullRequest)
	if err != nil {
		return nil, err
	}
	resp := ReassignReviewerResponse{
		NewReviewerID: newReviewer.ID(),
		PullRequest:   *pullRequest,
	}
	return &resp, nil
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

func (s *DefaultPullRequestDomainService) MarkAsMerged(pullRequestID v.ID) (*e.PullRequest, error) {
	pr, err := s.prRepo.FindByID(pullRequestID)
	if err != nil {
		return nil, err
	}
	if pr == nil {
		return nil, fmt.Errorf("no such pull request with id=%s", pullRequestID.String())
	}
	if pr.Status() == v.MERGED {
		return pr, nil
	}

	pr.MarkAsMerged()
	err = s.prRepo.Update(pr)
	if err != nil {
		return nil, err
	}
	return pr, nil
}
