package database

import (
	"context"
	"database/sql"
)

func NewUser(client *sql.DB) User {
	user := &userImpl{
		db: client,
	}
	return user
}

type UserAccount struct {
	ID       uint64
	Username string
	Password string
}

type User interface {
	GetUserByUsername(context.Context, string) (*UserAccount, error)
}

type userImpl struct {
	db *sql.DB
}

const sqlGetUserByUsername = `
SELECT id, username, passwd FROM user_account
WHERE username = ?
LIMIT 1
`

func (u *userImpl) GetUserByUsername(ctx context.Context, username string) (*UserAccount, error) {
	acc := &UserAccount{}
	if err := u.db.QueryRowContext(ctx, sqlGetUserByUsername, username).Scan(&acc.ID, &acc.Username, &acc.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoRecord
		}
		log.Errorf("failed to get user by username err: %v", err)
		return nil, ErrOpt
	}
	return acc, nil
}
