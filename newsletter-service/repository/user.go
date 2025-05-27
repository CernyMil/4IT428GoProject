package repository

import (
	"context"

	dbmodel "newsletter-management-api/repository/sql/model"
	"newsletter-management-api/repository/sql/query"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (r *UserRepository) ReadUser(ctx context.Context, userID string) (*dbmodel.User, error) {
	var user dbmodel.User
	if err := pgxscan.Get(
		ctx,
		r.pool,
		&user,
		query.ReadUser,
		pgx.NamedArgs{
			"id": userID,
		},
	); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) ListUser(ctx context.Context) ([]dbmodel.User, error) {
	var users []dbmodel.User
	if err := pgxscan.Select(
		ctx,
		r.pool,
		&users,
		query.ListUser,
	); err != nil {
		return nil, err
	}
	return users, nil
}
