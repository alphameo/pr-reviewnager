package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/alphameo/pr-reviewnager/internal/domain"
	db "github.com/alphameo/pr-reviewnager/internal/infra/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type TeamRepository struct {
	queries *db.Queries
	dbConn  *pgx.Conn
}

func NewTeamRepository(queries *db.Queries, databaseConnection *pgx.Conn) (*TeamRepository, error) {
	if queries == nil {
		return nil, errors.New("queries cannot be nil")
	}
	if databaseConnection == nil {
		return nil, errors.New("database connection cannot be nil")
	}

	return &TeamRepository{
		queries: queries,
		dbConn:  databaseConnection,
	}, nil
}

func (r *TeamRepository) Create(team *domain.Team) error {
	ctx := context.Background()
	tx, err := r.dbConn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	err = qtx.CreateTeam(ctx, db.CreateTeamParams{
		ID:   team.ID().Value(),
		Name: team.Name().Value(),
	})
	if err != nil {
		return err
	}

	for _, userID := range team.UserIDs() {
		userTeamID, err := qtx.GetTeamIDForUser(ctx, userID.Value())
		if err != nil {
			return err
		}
		if userTeamID != uuid.Nil {
			return fmt.Errorf("user with id=%s is already in team=%s", userID, userTeamID)
		}

		err = qtx.CreateTeamUser(ctx, db.CreateTeamUserParams{
			TeamID: team.ID().Value(),
			UserID: userID.Value(),
		})
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *TeamRepository) FindByID(id domain.ID) (*domain.Team, error) {
	ctx := context.Background()
	tx, err := r.dbConn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	teamID := id.Value()
	// TODO: rewrite with single query?
	dbTeam, err := qtx.GetTeam(ctx, teamID)
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	uIDs, err := qtx.GetUserIDsInTeam(ctx, teamID)
	if err != nil {
		return nil, err
	}
	userIDs := make([]domain.ID, 0, len(uIDs))
	for _, userID := range uIDs {
		userIDs = append(userIDs, domain.ExistingID(userID))
	}

	team := domain.ExistingTeam(
		domain.ExistingID(dbTeam.ID),
		domain.ExistingTeamName(dbTeam.Name),
		userIDs,
	)

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return team, nil
}

type TeamDTO struct {
	ID      domain.ID
	Name    string
	UserIDs []domain.ID
}

func (r *TeamRepository) FindAll() ([]*domain.Team, error) {
	ctx := context.Background()
	rows, err := r.queries.GetTeamsWithUsers(ctx)
	if err != nil {
		return nil, err
	}

	teamMap := make(map[uuid.UUID]*TeamDTO)

	for _, row := range rows {
		teamID := row.TeamID

		team, exists := teamMap[teamID]
		if !exists {
			team = &TeamDTO{
				ID:      domain.ExistingID(row.TeamID),
				Name:    row.TeamName,
				UserIDs: make([]domain.ID, 0),
			}
			teamMap[teamID] = team
		}

		if row.UserID.Valid {
			userID, err := domain.ParseID(row.UserID.String())
			if err != nil {
				return nil, err
			}
			team.UserIDs = append(team.UserIDs, userID)
		}
	}

	teams := make([]*domain.Team, 0)
	for _, teamDTO := range teamMap {
		team := domain.ExistingTeam(
			teamDTO.ID,
			domain.ExistingTeamName(teamDTO.Name),
			teamDTO.UserIDs,
		)
		teams = append(teams, team)
	}

	return teams, nil
}

func (r *TeamRepository) Update(team *domain.Team) error {
	ctx := context.Background()
	tx, err := r.dbConn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	err = qtx.UpdateTeam(ctx, db.UpdateTeamParams{
		ID:   team.ID().Value(),
		Name: team.Name().Value(),
	})
	if err != nil {
		return err
	}

	err = qtx.DeleteTeamUsersByTeamID(ctx, team.ID().Value())
	if err != nil {
		return err
	}

	for _, userID := range team.UserIDs() {
		userTeamID, err := qtx.GetTeamIDForUser(ctx, userID.Value())
		if err != nil {
			return err
		}
		if userTeamID != uuid.Nil {
			return fmt.Errorf("user with id=%s is already in team=%s", userID, userTeamID)
		}

		err = qtx.CreateTeamUser(
			ctx,
			db.CreateTeamUserParams{
				TeamID: team.ID().Value(),
				UserID: userID.Value(),
			})
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *TeamRepository) DeleteByID(id domain.ID) error {
	ctx := context.Background()

	err := r.queries.DeleteTeam(ctx, id.Value())
	if err != nil {
		return err
	}

	return nil
}

func (r *TeamRepository) FindByName(teamName string) (*domain.Team, error) {
	ctx := context.Background()
	tx, err := r.dbConn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	dbTeam, err := qtx.GetTeamByName(ctx, teamName)
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	uIDs, err := qtx.GetUserIDsInTeam(ctx, dbTeam.ID)
	if err != nil {
		return nil, err
	}
	userIDs := make([]domain.ID, 0, len(uIDs))
	for _, userID := range uIDs {
		userIDs = append(userIDs, domain.ExistingID(userID))
	}

	team := domain.ExistingTeam(
		domain.ExistingID(dbTeam.ID),
		domain.ExistingTeamName(dbTeam.Name),
		userIDs,
	)

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return team, nil
}

func (r *TeamRepository) CreateTeamAndModifyUsers(team *domain.Team, users []*domain.User) error {
	ctx := context.Background()
	tx, err := r.dbConn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	err = qtx.CreateTeam(ctx, db.CreateTeamParams{
		ID:   team.ID().Value(),
		Name: team.Name().Value(),
	})
	if err != nil {
		return err
	}

	for _, userID := range team.UserIDs() {
		userTeamID, err := qtx.GetTeamIDForUser(ctx, userID.Value())
		if err != nil {
			return err
		}
		if userTeamID != uuid.Nil {
			return fmt.Errorf("user with id=%s is already in team=%s", userID, userTeamID)
		}

		err = qtx.CreateTeamUser(ctx, db.CreateTeamUserParams{
			TeamID: team.ID().Value(),
			UserID: userID.Value(),
		})
		if err != nil {
			return err
		}
	}

	for _, user := range users {
		err = qtx.UpsertUser(ctx, db.UpsertUserParams{
			ID:     user.ID().Value(),
			Name:   user.Name().Value(),
			Active: user.Active(),
		})
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *TeamRepository) FindTeamByTeammateID(userID domain.ID) (*domain.Team, error) {
	ctx := context.Background()
	tx, err := r.dbConn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	dbTeam, err := qtx.GetTeamForUser(ctx, userID.Value())
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	uIDs, err := qtx.GetUserIDsInTeam(ctx, dbTeam.ID)
	if err != nil {
		return nil, err
	}
	userIDs := make([]domain.ID, 0, len(uIDs))
	for _, userID := range uIDs {
		userIDs = append(userIDs, domain.ExistingID(userID))
	}

	team := domain.ExistingTeam(
		domain.ExistingID(dbTeam.ID),
		domain.ExistingTeamName(dbTeam.Name),
		userIDs,
	)

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return team, nil
}

func (r *TeamRepository) FindActiveUsersByTeamID(teamID domain.ID) ([]*domain.User, error) {
	ctx := context.Background()

	users, err := r.queries.GetActiveUsersInTeam(ctx, teamID.Value())
	if err != nil {
		return nil, err
	}

	entities := make([]*domain.User, len(users))
	for i, user := range users {
		entities[i] = domain.ExistingUser(
			domain.ExistingID(user.ID),
			domain.ExistingUserName(user.Name),
			user.Active,
		)
	}

	return entities, nil
}
