package repo

import (
	"context"
	"fmt"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/pkg/postgres"
)

const UsersDBName = "users"

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

func (ur *UserRepo) CreateUser(ctx context.Context, user *entity.User) (int, error) {
	sql, args, err := ur.Builder.Insert(UsersDBName).Columns("name", "email", "age", "password").
		Values(user.Name, user.Email, user.Age, user.Password).
		Suffix("returning id").ToSql()
	if err != nil {
		return 0, err
	}

	var insertedID int

	err = ur.Pool.QueryRow(ctx, sql, args...).Scan(&insertedID)
	if err != nil {
		return 0, err
	}

	return insertedID, nil
}

func (ur *UserRepo) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	sql, args, err := ur.Builder.Select("id", "name", "email", "age", "password").
		From("users").
		Where("email = $1", email).
		ToSql()

	if err != nil {
		return nil, err
	}

	user := new(entity.User)
	err = ur.Pool.QueryRow(ctx, sql, args...).Scan(&user.Id, &user.Name, &user.Email, &user.Age, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}
