package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/alphameo/pr-reviewnager/internal/domain/dto"
	e "github.com/alphameo/pr-reviewnager/internal/domain/entities"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
	db "github.com/alphameo/pr-reviewnager/internal/infrastructure/db/sqlc"
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

func (r *TeamRepository) Create(team *e.Team) error {
	ctx := context.Background()
	tx, err := r.dbConn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	err = qtx.CreateTeam(ctx, db.CreateTeamParams{
		ID:   uuid.UUID(team.ID()),
		Name: team.Name(),
	})
	if err != nil {
		return err
	}

	for _, userID := range team.UserIDs() {
		userTeamID, err := qtx.GetTeamIDForUser(ctx, uuid.UUID(userID))
		if err != nil {
			return err
		}
		if userTeamID != uuid.Nil {
			return fmt.Errorf("user with id=%s is already in team=%s", userID, userTeamID)
		}

		err = qtx.CreateTeamUser(ctx, db.CreateTeamUserParams{
			TeamID: uuid.UUID(team.ID()),
			UserID: uuid.UUID(userID),
		})
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// TODO: rewrite with single querry?
func (r *TeamRepository) FindByID(id v.ID) (*dto.TeamDTO, error) {
	ctx := context.Background()
	tx, err := r.dbConn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	teamID := uuid.UUID(id)
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
	userIDs := make([]v.ID, 0, len(uIDs))
	for _, userID := range uIDs {
		userIDs = append(userIDs, v.ID(userID))
	}

	team := dto.TeamDTO{
		ID:      v.ID(dbTeam.ID),
		Name:    dbTeam.Name,
		UserIDs: userIDs,
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &team, nil
}

// TODO: rewrite with single querry?
func (r *TeamRepository) FindAll() ([]*dto.TeamDTO, error) {
	ctx := context.Background()
	tx, err := r.dbConn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	rows, err := qtx.GetTeamsWithUsers(ctx)
	if err != nil {
		return nil, err
	}

	teamMap := make(map[uuid.UUID]*dto.TeamDTO)
	teams := make([]*dto.TeamDTO, 0)

	for _, row := range rows {
		teamID := uuid.UUID(row.TeamID)

		team, exists := teamMap[teamID]
		if !exists {
			team = &dto.TeamDTO{
				ID:      v.ID(row.TeamID),
				Name:    row.TeamName,
				UserIDs: make([]v.ID, 0),
			}
			teamMap[teamID] = team
			teams = append(teams, team)
		}

		if row.UserID.Valid {
			userID, err := v.NewIDFromString(row.UserID.String())
			if err != nil {
				return nil, err
			}
			team.UserIDs = append(team.UserIDs, userID)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return teams, nil
}

func (r *TeamRepository) Update(team *e.Team) error {
	ctx := context.Background()
	tx, err := r.dbConn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	err = qtx.UpdateTeam(ctx, db.UpdateTeamParams{
		ID:   uuid.UUID(team.ID()),
		Name: team.Name(),
	})
	if err != nil {
		return err
	}

	err = qtx.DeleteTeamUsersByTeamID(ctx, uuid.UUID(team.ID()))
	if err != nil {
		return err
	}

	for _, userID := range team.UserIDs() {
		userTeamID, err := qtx.GetTeamIDForUser(ctx, uuid.UUID(userID))
		if err != nil {
			return err
		}
		if userTeamID != uuid.Nil {
			return fmt.Errorf("user with id=%s is already in team=%s", userID, userTeamID)
		}

		err = qtx.CreateTeamUser(
			ctx,
			db.CreateTeamUserParams{
				TeamID: uuid.UUID(team.ID()),
				UserID: uuid.UUID(userID),
			})
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *TeamRepository) DeleteByID(id v.ID) error {
	ctx := context.Background()

	err := r.queries.DeleteTeam(ctx, uuid.UUID(id))
	if err != nil {
		return err
	}

	return nil
}

func (r *TeamRepository) FindByName(teamName string) (*dto.TeamDTO, error) {
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
	userIDs := make([]v.ID, 0, len(uIDs))
	for _, userID := range uIDs {
		userIDs = append(userIDs, v.ID(userID))
	}

	team := dto.TeamDTO{
		ID:      v.ID(dbTeam.ID),
		Name:    dbTeam.Name,
		UserIDs: userIDs,
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &team, nil
}

func (r *TeamRepository) CreateTeamAndModifyUsers(team *e.Team, users []*e.User) error {
	ctx := context.Background()
	tx, err := r.dbConn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	err = qtx.CreateTeam(ctx, db.CreateTeamParams{
		ID:   uuid.UUID(team.ID()),
		Name: team.Name(),
	})
	if err != nil {
		return err
	}

	for _, userID := range team.UserIDs() {
		userTeamID, err := qtx.GetTeamIDForUser(ctx, uuid.UUID(userID))
		if err != nil {
			return err
		}
		if userTeamID != uuid.Nil {
			return fmt.Errorf("user with id=%s is already in team=%s", userID, userTeamID)
		}

		err = qtx.CreateTeamUser(ctx, db.CreateTeamUserParams{
			TeamID: uuid.UUID(team.ID()),
			UserID: uuid.UUID(userID),
		})
		if err != nil {
			return err
		}
	}

	for _, user := range users {
		err = qtx.UpsertUser(ctx, db.UpsertUserParams{
			ID:     uuid.UUID(user.ID()),
			Name:   user.Name(),
			Active: user.Active(),
		})
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *TeamRepository) FindTeamByTeammateID(userID v.ID) (*dto.TeamDTO, error) {
	ctx := context.Background()
	tx, err := r.dbConn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	dbTeam, err := qtx.GetTeamForUser(ctx, uuid.UUID(userID))
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	uIDs, err := qtx.GetUserIDsInTeam(ctx, dbTeam.ID)
	if err != nil {
		return nil, err
	}
	userIDs := make([]v.ID, 0, len(uIDs))
	for _, userID := range uIDs {
		userIDs = append(userIDs, v.ID(userID))
	}

	team := dto.TeamDTO{
		ID:      v.ID(dbTeam.ID),
		Name:    dbTeam.Name,
		UserIDs: userIDs,
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &team, nil
}

func (r *TeamRepository) FindActiveUsersByTeamID(teamID v.ID) ([]*dto.UserDTO, error) {
	ctx := context.Background()

	users, err := r.queries.GetActiveUsersInTeam(ctx, uuid.UUID(teamID))
	if err != nil {
		return nil, err
	}

	entities := make([]*dto.UserDTO, len(users))
	for i, user := range users {
		entities[i] = &dto.UserDTO{
			ID:     v.ID(user.ID),
			Name:   user.Name,
			Active: user.Active,
		}
	}

	return entities, nil
}
