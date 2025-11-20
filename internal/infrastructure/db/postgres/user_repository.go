package postgres

import (
	"context"
	"errors"

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

	r := UserRepository{queries: queries}
	return &r, nil
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

func (r *UserRepository) FindByID(id v.ID) (*e.User, error) {
	ctx := context.Background()

	user, err := r.queries.GetUser(ctx, uuid.UUID(id))
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	entity := e.NewExistingUser(v.ID(user.ID), user.Name, user.Active)

	return entity, nil
}

func (r *UserRepository) FindAll() ([]*e.User, error) {
	ctx := context.Background()

	users, err := r.queries.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	entities := make([]*e.User, len(users))
	for i, user := range users {
		u := e.NewExistingUser(v.ID(user.ID), user.Name, user.Active)
		entities[i] = u
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
