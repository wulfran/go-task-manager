package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"task-manager/internal/db"
	"task-manager/internal/helpers"
	"task-manager/internal/models"
	"time"
)

type UserRepository interface {
	CreateUser(ctx context.Context, r models.CreateUserPayload) error
	CheckIfEmailExists(email string) (bool, error)
	GetUserData(p models.LoginPayload) (models.User, error)
}

type userRepository struct {
	db db.DB
}

func NewUserRepository(d db.DB) UserRepository {
	return &userRepository{
		db: d,
	}
}

func (u userRepository) CreateUser(ctx context.Context, r models.CreateUserPayload) error {
	q, err := db.GetQuery("user/InsertUser.sql")
	if err != nil {
		return fmt.Errorf("CreateUser: error while reading query: %v", err)
	}

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("CreateUser: failed to begin tx: %v", err)
	}

	_, err = tx.ExecContext(
		ctx,
		q,
		r.Name,
		r.Email,
		r.Password,
		time.Now().Format(time.RFC3339),
	)

	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("CreateUser: failed to insert a new user: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("CreateUser: failed to commit tx: %v", err)
	}

	return nil
}
func (u userRepository) CheckIfEmailExists(email string) (bool, error) {
	q, err := db.GetQuery("user/EmailExistsWithinUsers.sql")
	if err != nil {
		return false, fmt.Errorf("CheckIfEmailExists: error while reading query: %v", err)
	}

	var exists bool

	err = u.db.QueryRow(q, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("CheckIfEmailExists: failed to execute query: %v", err)
	}

	return exists, nil
}
func (u userRepository) GetUserData(p models.LoginPayload) (models.User, error) {
	q, err := db.GetQuery("user/LoginUser.sql")
	if err != nil {
		return models.User{}, fmt.Errorf("GetUserData: error while reading query: %v", err)
	}
	var uData models.User

	err = u.db.QueryRow(q, p.Email).Scan(&uData.ID, &uData.Name, &uData.Email, &uData.Password, &uData.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, fmt.Errorf("GetUserData: no entries found")
	}

	if err != nil {
		return models.User{}, fmt.Errorf("GetUserData: failed to execute query: %v", err)
	}

	if !helpers.ValidatePassword(p.Password, uData.Password) {
		return models.User{}, fmt.Errorf("incorrect credentials")
	}

	return uData, nil
}
