package postgres

import (
	"context"
	"errors"

	"github.com/alphameo/pr-reviewnager/internal/domain/dto"
	e "github.com/alphameo/pr-reviewnager/internal/domain/entities"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
	db "github.com/alphameo/pr-reviewnager/internal/infrastructure/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	queries *db.Queries
}

func NewUserRepository(queries *db.Queries) (*UserRepository, error) {
	if queries == nil {
		return nil, errors.New("queries cannot be nil")
	}

	return &UserRepository{queries: queries}, nil
}

func (r *UserRepository) Create(user *e.User) error {
	ctx := context.Background()

	if user == nil {
		return errors.New("user cannot be nil")
	}

	err := r.queries.CreateUser(ctx, db.CreateUserParams{
		ID:     uuid.UUID(user.ID()),
		Name:   user.Name(),
		Active: user.Active(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) FindByID(id v.ID) (*dto.User, error) {
	ctx := context.Background()

	user, err := r.queries.GetUser(ctx, uuid.UUID(id))
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &dto.User{
		ID:     v.ID(user.ID),
		Name:   user.Name,
		Active: user.Active,
	}, nil
}

func (r *UserRepository) FindAll() ([]*dto.User, error) {
	ctx := context.Background()

	users, err := r.queries.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	entities := make([]*dto.User, len(users))
	for i, user := range users {
		entities[i] = &dto.User{
			ID:     v.ID(user.ID),
			Name:   user.Name,
			Active: user.Active,
		}
	}

	return entities, nil
}

func (r *UserRepository) Update(user *e.User) error {
	ctx := context.Background()

	if user == nil {
		return errors.New("user cannot be nil")
	}

	err := r.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:     uuid.UUID(user.ID()),
		Name:   user.Name(),
		Active: user.Active(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) DeleteByID(id v.ID) error {
	ctx := context.Background()

	err := r.queries.DeleteUser(ctx, uuid.UUID(id))
	if err != nil {
		return err
	}

	return nil
}
