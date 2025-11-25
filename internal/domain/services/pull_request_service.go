// Package services provides domain services for domain model
package services

import (
	"errors"
	"fmt"
	"math/rand"
	"slices"

	"github.com/alphameo/pr-reviewnager/internal/domain/dto"
	e "github.com/alphameo/pr-reviewnager/internal/domain/entities"
	r "github.com/alphameo/pr-reviewnager/internal/domain/repositories"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
)

type PullRequestDomainService interface {
	// CreateWithReviewers() creates a new pull request and automatically assigns
	// 2 reviewers by randomly selecting from the author's team.
	CreateWithReviewers(pullRequest *e.PullRequest) (*e.PullRequest, error)

	// ReassignReviewer() unassign user-reviewer with given id and assigns another from his team, excluding
	// him and pr author. After, method returns id of new user-reviewer and pull request
	ReassignReviewer(userID v.ID, pullRequestID v.ID) (*ReassignReviewerResponse, error)

	// MarkAsMerged() idempotently marks pull request as merged and sets time of marking
	MarkAsMerged(pullRequestID v.ID) (*e.PullRequest, error)
}

type DefaultPullRequestDomainService struct {
	userRepo r.UserRepository
	teamRepo r.TeamRepository
	prRepo   r.PullRequestRepository
}

var (
	ErrAuthorNotFound     error = errors.New("author not found")
	ErrTeamNotFound       error = errors.New("team not found")
	ErrPRAlreadyExists    error = errors.New("pull request already exists")
	ErrPRNotFound         error = errors.New("pull request not found")
	ErrUserNotFound       error = errors.New("user not found")
	ErrUserNotReviewer    error = errors.New("user is not a reviewer")
	ErrNoReviewCandidates error = errors.New("no users ready to review")
)

func NewDefaultPullRequestDomainService(
	userRepository r.UserRepository,
	pullRequestRepository r.PullRequestRepository,
	teamRepository r.TeamRepository,
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

	return &DefaultPullRequestDomainService{
		userRepo: userRepository,
		prRepo:   pullRequestRepository,
		teamRepo: teamRepository,
	}, nil
}

func (s *DefaultPullRequestDomainService) CreateWithReviewers(pullRequest *e.PullRequest) (*e.PullRequest, error) {
	prDTO, err := s.prRepo.FindByID(pullRequest.ID())
	if err != nil {
		return nil, err
	}
	if prDTO != nil {
		return nil, ErrPRAlreadyExists
	}

	authorID := pullRequest.AuthorID()

	author, err := s.userRepo.FindByID(authorID)
	if err != nil {
		return nil, err
	}
	if author == nil {
		return nil, ErrAuthorNotFound
	}

	team, err := s.teamRepo.FindTeamByTeammateID(authorID)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, ErrTeamNotFound
	}

	availableUsers, err := s.teamRepo.FindActiveUsersByTeamID(team.ID)
	if err != nil {
		return nil, err
	}

	reviewers := chooseRandomUsers(availableUsers, e.MaxReviewersCount, authorID)
	for _, u := range reviewers {
		pullRequest.AssignReviewer(u.ID)
	}

	err = s.prRepo.Create(pullRequest)
	if err != nil {
		return nil, err
	}

	return pullRequest, nil
}

type ReassignReviewerResponse struct {
	NewReviewerID v.ID
	PullRequest   e.PullRequest
}

func (s *DefaultPullRequestDomainService) ReassignReviewer(userID v.ID, pullRequestID v.ID) (*ReassignReviewerResponse, error) {
	prDTO, err := s.prRepo.FindByID(pullRequestID)
	if err != nil {
		return nil, err
	}
	if prDTO == nil {
		return nil, ErrPRNotFound
	}

	pr, err := e.NewExistingPullRequest(prDTO)
	if err != nil {
		return nil, err
	}
	if pr.Status() == v.MERGED {
		return nil, e.ErrPRAlreadyMerged
	}

	reviewerIDs := pr.ReviewerIDs()
	reviewerIdx := slices.Index(reviewerIDs, userID)
	if reviewerIdx == -1 {
		return nil, fmt.Errorf("cannot reassign reviewer with id=%s: %w", userID.String(), ErrUserNotReviewer)
	}

	authorID := pr.AuthorID()
	team, err := s.teamRepo.FindTeamByTeammateID(authorID)
	if err != nil {
		return nil, err
	}

	availableUsers, err := s.teamRepo.FindActiveUsersByTeamID(team.ID)
	if err != nil {
		return nil, err
	}
	if len(availableUsers) == 0 {
		return nil, ErrNoReviewCandidates
	}

	exceptionalReviewerIDs := pr.ReviewerIDs()
	exceptionalReviewerIDs = append(exceptionalReviewerIDs, authorID)
	newReviewer := chooseRandomUser(availableUsers, exceptionalReviewerIDs...)

	pr.UnassignReviewer(userID)
	pr.AssignReviewer(newReviewer.ID)
	err = s.prRepo.Update(pr)
	if err != nil {
		return nil, err
	}

	return &ReassignReviewerResponse{
		NewReviewerID: newReviewer.ID,
		PullRequest:   *pr,
	}, nil
}

func chooseRandomUser(availableUsers []*dto.UserDTO, except ...v.ID) *dto.UserDTO {
	for i, u := range availableUsers {
		for _, exceptionalU := range except {
			if u.ID == exceptionalU {
				availableUsers = slices.Delete(availableUsers, i, i+1)
				continue
			}
		}
	}
	idx := rand.Intn(len(availableUsers))

	return availableUsers[idx]
}

func chooseRandomUsers(availableUsers []*dto.UserDTO, maxCount int, except ...v.ID) []*dto.UserDTO {
	for i, u := range availableUsers {
		for _, exceptionalU := range except {
			if u.ID == exceptionalU {
				availableUsers = slices.Delete(availableUsers, i, i+1)
				continue
			}
		}
	}

	reviewers := make([]*dto.UserDTO, 0, maxCount)
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
	prDTO, err := s.prRepo.FindByID(pullRequestID)
	if err != nil {
		return nil, err
	}
	if prDTO == nil {
		return nil, ErrPRNotFound
	}

	pr, err := e.NewExistingPullRequest(prDTO)
	if err != nil {
		return nil, err
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
