package api

import (
	"errors"
	"net/http"

	app "github.com/alphameo/pr-reviewnager/internal/application"
	"github.com/alphameo/pr-reviewnager/internal/domain"
	"github.com/labstack/echo/v4"
)

type Server struct {
	teamService app.TeamService
	userService app.UserService
	prService   app.PullRequestService
}

func NewServer(teamService app.TeamService, userService app.UserService, pullRequestService app.PullRequestService) (*Server, error) {
	if teamService == nil {
		return nil, errors.New("teamService cannot be nil")
	}
	if userService == nil {
		return nil, errors.New("userService cannot be nil")
	}
	if pullRequestService == nil {
		return nil, errors.New("pullRequestService cannot be nil")
	}

	return &Server{
		teamService: teamService,
		userService: userService,
		prService:   pullRequestService,
	}, nil
}

func (s Server) PostPullRequestCreate(ctx echo.Context) error {
	var input PostPullRequestCreateJSONRequestBody
	if err := ctx.Bind(&input); err != nil {
		return err
	}

	prID, err := domain.NewIDFromString(input.PullRequestId)
	if err != nil {
		return err
	}
	authorID, err := domain.NewIDFromString(input.AuthorId)
	if err != nil {
		return err
	}

	req := domain.PullRequestDTO{
		ID:          prID,
		Title:       input.PullRequestName,
		AuthorID:    authorID,
		Status:      "OPEN",
		ReviewerIDs: nil,
	}

	createdPR, err := s.prService.CreatePullRequest(&req)
	if err != nil {
		return mapAppErrorToEchoResponse(ctx, err)
	}

	return ctx.JSON(http.StatusCreated, map[string]PullRequest{
		"pr": ToAPIPullRequest(*createdPR),
	})
}

func (s *Server) PostPullRequestMerge(ctx echo.Context) error {
	var input PostPullRequestMergeJSONRequestBody
	if err := ctx.Bind(&input); err != nil {
		return err
	}

	prID, err := domain.NewIDFromString(input.PullRequestId)
	if err != nil {
		return err
	}

	dtoPR, err := s.prService.MarkAsMerged(prID)
	if err != nil {
		return mapAppErrorToEchoResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, map[string]PullRequest{
		"pr": ToAPIPullRequest(*dtoPR),
	})
}

func (s *Server) PostPullRequestReassign(ctx echo.Context) error {
	var input PostPullRequestReassignJSONRequestBody
	if err := ctx.Bind(&input); err != nil {
		return err
	}

	prID, err := domain.NewIDFromString(input.PullRequestId)
	if err != nil {
		return err
	}
	oldID, err := domain.NewIDFromString(input.OldUserId)
	if err != nil {
		return err
	}

	resp, err := s.prService.ReassignReviewer(prID, oldID)
	if err != nil {
		return mapAppErrorToEchoResponse(ctx, err)
	}
	updatedPR := resp.PullRequest
	replacedBy := resp.NewReviewerUserID

	return ctx.JSON(http.StatusOK, map[string]any{
		"pr":          ToAPIPullRequest(*updatedPR),
		"replaced_by": replacedBy.String(),
	})
}

func (s *Server) PostTeamAdd(ctx echo.Context) error {
	var team Team
	if err := ctx.Bind(&team); err != nil {
		return err
	}

	teamDTO := FromAPITeam(team)

	err := s.teamService.CreateTeamWithUsers(&teamDTO)
	if err != nil {
		return mapAppErrorToEchoResponse(ctx, err)
	}

	return ctx.JSON(http.StatusCreated, map[string]Team{
		"team": team,
	})
}

func (s *Server) GetTeamGet(ctx echo.Context, params GetTeamGetParams) error {
	dtoTeam, err := s.teamService.FindTeamByName(params.TeamName)
	if err != nil {
		return mapAppErrorToEchoResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, ToAPITeam(*dtoTeam))
}

func (s *Server) PostUsersSetIsActive(ctx echo.Context) error {
	var input PostUsersSetIsActiveJSONRequestBody
	if err := ctx.Bind(&input); err != nil {
		return err
	}

	userID, err := domain.NewIDFromString(input.UserId)
	if err != nil {
		return err
	}

	updated, err := s.teamService.SetUserActiveByID(userID, input.IsActive)
	if err != nil {
		return mapAppErrorToEchoResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, map[string]User{
		"user": ToAPIUser(*updated),
	})
}

func (s *Server) GetUsersGetReview(ctx echo.Context, params GetUsersGetReviewParams) error {
	userID, err := domain.NewIDFromString(params.UserId)
	if err != nil {
		return err
	}

	list, err := s.prService.FindPullRequestsByReviewer(userID)
	if err != nil {
		return mapAppErrorToEchoResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, map[string]any{
		"user_id":       params.UserId,
		"pull_requests": ToAPIPullRequestShortList(list),
	})
}

func mapAppErrorToEchoResponse(ctx echo.Context, err error) error {
	switch {
	case errors.Is(err, app.ErrTeamExists):
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error: struct {
				Code    ErrorResponseErrorCode `json:"code"`
				Message string                 `json:"message"`
			}{
				Code:    TEAMEXISTS,
				Message: "team_name already exists",
			},
		})

	case errors.Is(err, app.ErrPRExists):
		return ctx.JSON(http.StatusConflict, ErrorResponse{
			Error: struct {
				Code    ErrorResponseErrorCode `json:"code"`
				Message string                 `json:"message"`
			}{
				Code:    PREXISTS,
				Message: "PR id already exists",
			},
		})

	case errors.Is(err, app.ErrPRAlreadyMerged):
		return ctx.JSON(http.StatusConflict, ErrorResponse{
			Error: struct {
				Code    ErrorResponseErrorCode `json:"code"`
				Message string                 `json:"message"`
			}{
				Code:    PRMERGED,
				Message: "cannot reassign on merged PR",
			},
		})

	case errors.Is(err, app.ErrNotAssigned):
		return ctx.JSON(http.StatusConflict, ErrorResponse{
			Error: struct {
				Code    ErrorResponseErrorCode `json:"code"`
				Message string                 `json:"message"`
			}{
				Code:    NOTASSIGNED,
				Message: "reviewer is not assigned to this PR",
			},
		})

	case errors.Is(err, app.ErrNoCandidate):
		return ctx.JSON(http.StatusConflict, ErrorResponse{
			Error: struct {
				Code    ErrorResponseErrorCode `json:"code"`
				Message string                 `json:"message"`
			}{
				Code:    NOCANDIDATE,
				Message: "no active replacement candidate in team",
			},
		})

	case errors.Is(err, app.ErrNotFound):
		return ctx.JSON(http.StatusNotFound, ErrorResponse{
			Error: struct {
				Code    ErrorResponseErrorCode `json:"code"`
				Message string                 `json:"message"`
			}{
				Code:    NOTFOUND,
				Message: "resource not found",
			},
		})
	}

	return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: struct {
			Code    ErrorResponseErrorCode `json:"code"`
			Message string                 `json:"message"`
		}{
			Code:    "INTERNAL",
			Message: "internal server error",
		},
	})
}
