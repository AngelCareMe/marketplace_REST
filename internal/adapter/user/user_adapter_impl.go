package adapter

import (
	"context"
	"database/sql"
	"fmt"
	"marketplace/internal/entity"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type UserAdapter struct {
	db     *pgxpool.Pool
	logger *logrus.Logger
}

func NewUserAdaper(db *pgxpool.Pool, logger *logrus.Logger) *UserAdapter {
	return &UserAdapter{
		db:     db,
		logger: logger,
	}
}

func (a *UserAdapter) Create(ctx context.Context, user *entity.User) error {
	query, args, err := squirrel.
		Insert("users").
		Columns("id", "username", "hashed_password", "created_at").
		Values(user.ID, user.Username, user.HashedPassword, user.CreatedAt).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		a.logger.WithError(err).Error("Failed to build create user query")
		return err
	}

	_, err = a.db.Exec(ctx, query, args...)
	if err != nil {
		a.logger.WithError(err).Error("Failed to create user")
		return err
	}

	a.logger.WithFields(logrus.Fields{
		"user_id":  user.ID,
		"username": user.Username,
	}).Info("User created in database")
	return nil
}

func (a *UserAdapter) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	query, args, err := squirrel.Select("id", "username", "hashed_password", "created_at").
		From("users").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		a.logger.WithError(err).Error("Failed to build get user by ID query")
		return nil, err
	}

	var user entity.User
	err = a.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Username, &user.HashedPassword, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		a.logger.WithError(err).Error("Failed to get user by ID")
		return nil, err
	}

	return &user, nil
}

func (a *UserAdapter) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	query, args, err := squirrel.Select("id", "username", "hashed_password", "created_at").
		From("users").
		Where(squirrel.Eq{"username": username}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		a.logger.WithError(err).Error("Failed to build get user by username query")
		return nil, err
	}

	var user entity.User
	err = a.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Username, &user.HashedPassword, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		a.logger.WithError(err).Error("Failed to get user by username")
		return nil, err
	}

	return &user, nil
}

func (a *UserAdapter) Update(ctx context.Context, user *entity.User) error {
	query, args, err := squirrel.Update("users").
		Set("username", user.Username).
		Set("hashed_password", user.HashedPassword).
		Where(squirrel.Eq{"id": user.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		a.logger.WithError(err).Error("Failed to build update user query")
		return err
	}

	result, err := a.db.Exec(ctx, query, args...)
	if err != nil {
		a.logger.WithError(err).Error("Failed to update user")
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	a.logger.WithFields(logrus.Fields{
		"user_id": user.ID,
	}).Info("User updated in database")
	return nil
}

func (a *UserAdapter) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := squirrel.Delete("users").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		a.logger.WithError(err).Error("Failed to build delete user query")
		return err
	}

	result, err := a.db.Exec(ctx, query, args...)
	if err != nil {
		a.logger.WithError(err).Error("Failed to delete user")
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	a.logger.WithFields(logrus.Fields{
		"user_id": id,
	}).Info("User deleted from database")
	return nil
}
