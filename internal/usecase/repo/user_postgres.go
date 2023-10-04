package repo

import (
	"context"
	"fmt"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/pkg/postgres"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (ur *UserRepo) GetUsers(ctx context.Context) ([]*entity.User, error) {
	sql, _, err := ur.Builder.
		Select("id, name, email, age").
		From("users").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("users repo - get users - r.Builder: %w", err)
	}

	rows, err := ur.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("users repo - get users -- r.Pool.Query: %w", err)
	}
	defer rows.Close()

	entities := make([]*entity.User, 0, _defaultEntityCap)

	for rows.Next() {
		user := new(entity.User)

		err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Age)
		if err != nil {
			return nil, fmt.Errorf("users repo - get users - r.Pool.Query- rows.Scan: %w", err)
		}

		entities = append(entities, user)
	}

	return entities, nil
}
