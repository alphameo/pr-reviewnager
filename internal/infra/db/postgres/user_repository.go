package postgres

import (
	"context"
	"errors"

	"github.com/alphameo/pr-reviewnager/internal/domain"
	db "github.com/alphameo/pr-reviewnager/internal/infra/db/sqlc"
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

func (r *UserRepository) Create(user *domain.User) error {
	ctx := context.Background()

	if user == nil {
		return errors.New("user cannot be nil")
	}

	err := r.queries.CreateUser(ctx, db.CreateUserParams{
		ID:     user.ID().Value(),
		Name:   user.Name(),
		Active: user.Active(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) FindByID(id domain.ID) (*domain.User, error) {
	ctx := context.Background()

	user, err := r.queries.GetUser(ctx, id.Value())
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return domain.ExistingUser(
		domain.ExistingID(user.ID),
		user.Name,
		user.Active,
	), nil
}

func (r *UserRepository) FindAll() ([]*domain.User, error) {
	ctx := context.Background()

	users, err := r.queries.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	entities := make([]*domain.User, len(users))
	for i, user := range users {
		entities[i] = domain.ExistingUser(
			domain.ExistingID(user.ID),
			user.Name,
			user.Active,
		)
	}

	return entities, nil
}

func (r *UserRepository) Update(user *domain.User) error {
	ctx := context.Background()

	if user == nil {
		return errors.New("user cannot be nil")
	}

	err := r.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:     user.ID().Value(),
		Name:   user.Name(),
		Active: user.Active(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) DeleteByID(id domain.ID) error {
	ctx := context.Background()

	err := r.queries.DeleteUser(ctx, id.Value())
	if err != nil {
		return err
	}

	return nil
}
