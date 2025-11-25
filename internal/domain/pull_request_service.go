package domain

import (
	"errors"
	"fmt"
	"math/rand"
	"slices"
)

type PullRequestDomainService interface {
	// CreateWithReviewers() creates a new pull request and automatically assigns
	// 2 reviewers by randomly selecting from the author's team.
	CreateWithReviewers(pullRequest *PullRequest) (*PullRequest, error)

	// ReassignReviewer() unassign user-reviewer with given id and assigns another from his team, excluding
	// him and pr author. After, method returns id of new user-reviewer and pull request
	ReassignReviewer(userID ID, pullRequestID ID) (*ReassignReviewerResponse, error)

	// MarkAsMerged() idempotently marks pull request as merged and sets time of marking
	MarkAsMerged(pullRequestID ID) (*PullRequest, error)
}

type DefaultPullRequestDomainService struct {
	userRepo UserRepository
	teamRepo TeamRepository
	prRepo   PullRequestRepository
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
	userRepository UserRepository,
	pullRequestRepository PullRequestRepository,
	teamRepository TeamRepository,
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

func (s *DefaultPullRequestDomainService) CreateWithReviewers(pullRequest *PullRequest) (*PullRequest, error) {
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

	availableUsers, err := s.teamRepo.FindActiveUsersByTeamID(team.ID())
	if err != nil {
		return nil, err
	}

	reviewers := chooseRandomUsers(availableUsers, MaxReviewersCount, authorID)
	for _, u := range reviewers {
		pullRequest.AssignReviewer(u.ID())
	}

	err = s.prRepo.Create(pullRequest)
	if err != nil {
		return nil, err
	}

	return pullRequest, nil
}

type ReassignReviewerResponse struct {
	NewReviewerID ID
	PullRequest   PullRequest
}

func (s *DefaultPullRequestDomainService) ReassignReviewer(userID ID, pullRequestID ID) (*ReassignReviewerResponse, error) {
	pr, err := s.prRepo.FindByID(pullRequestID)
	if err != nil {
		return nil, err
	}
	if pr == nil {
		return nil, ErrPRNotFound
	}

	if pr.Status() == PRMerged {
		return nil, ErrPRAlreadyMerged
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

	availableUsers, err := s.teamRepo.FindActiveUsersByTeamID(team.ID())
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
	pr.AssignReviewer(newReviewer.ID())
	err = s.prRepo.Update(pr)
	if err != nil {
		return nil, err
	}

	return &ReassignReviewerResponse{
		NewReviewerID: newReviewer.ID(),
		PullRequest:   *pr,
	}, nil
}

func chooseRandomUser(availableUsers []*User, except ...ID) *User {
	for i, u := range availableUsers {
		for _, exceptional := range except {
			if u.ID() == exceptional {
				availableUsers = slices.Delete(availableUsers, i, i+1)
				continue
			}
		}
	}
	idx := rand.Intn(len(availableUsers))

	return availableUsers[idx]
}

func chooseRandomUsers(availableUsers []*User, maxCount int, except ...ID) []*User {
	for i, u := range availableUsers {
		for _, exceptional := range except {
			if u.ID() == exceptional {
				availableUsers = slices.Delete(availableUsers, i, i+1)
				continue
			}
		}
	}

	reviewers := make([]*User, 0, maxCount)
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

func (s *DefaultPullRequestDomainService) MarkAsMerged(pullRequestID ID) (*PullRequest, error) {
	pr, err := s.prRepo.FindByID(pullRequestID)
	if err != nil {
		return nil, err
	}
	if pr == nil {
		return nil, ErrPRNotFound
	}

	if pr.Status() == PRMerged {
		return pr, nil
	}

	pr.MarkAsMerged()
	err = s.prRepo.Update(pr)
	if err != nil {
		return nil, err
	}

	return pr, nil
}
