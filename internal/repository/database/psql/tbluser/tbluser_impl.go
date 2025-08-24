package tbluser

import (
	"backend-assignment/internal/models"
	"backend-assignment/internal/repository/database"
	"backend-assignment/logs"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type TblUserImpl struct {
	DB *sql.DB
}

func New(DB *sql.DB) database.TblUser {
	return &TblUserImpl{
		DB: DB,
	}
}

// Register implements database.TblUser.
func (repository *TblUserImpl) Register(ctx context.Context, request models.RequestUser) error {
	location, _ := time.LoadLocation("Asia/Jakarta")
	date := time.Now().In(location)

	tx, err := repository.DB.Begin()
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
		}
	}()

	sql := `INSERT INTO users (name, created_at) VALUES ($1, $2)`

	_, err = tx.ExecContext(ctx, sql, request.Name, date)
	if err != nil {
		logs.Error(err.Error())
		return err
	}

	err = tx.Commit()
	if err != nil {
		logs.Error(err.Error())
		return err
	}

	return nil
}

// GetUser implements database.TblUser.
func (repository *TblUserImpl) GetUser(ctx context.Context, userId int) (*models.User, error) {
	query := `SELECT name FROM users WHERE id = $1`

	var userDetail models.User
	err := repository.DB.QueryRowContext(ctx, query, userId).Scan(
		&userDetail.Name,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get order detail: %w", err)
	}

	return &userDetail, nil
}
