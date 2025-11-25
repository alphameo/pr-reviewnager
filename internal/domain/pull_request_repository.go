package domain


type PullRequestRepository interface {
	Repository[PullRequest, PullRequestDTO, ID]
	FindPullRequestsByReviewer(userID ID) ([]*PullRequestDTO, error)
}
