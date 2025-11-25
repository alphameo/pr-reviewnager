package domain


type PullRequestRepository interface {
	Repository[PullRequest, ID]
	FindPullRequestsByReviewer(userID ID) ([]*PullRequest, error)
}
