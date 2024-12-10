package db

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `json:"-" bun:"table:user"`
	ID            int    `bun:"id,autoincrement" json:"id"`
	UserName      string `bun:"username,pk" binding:"required" json:"username"`
	Password      string `bun:"password" binding:"required" json:"password"`
	FirstName     string `bun:"first_name" binding:"required" json:"first_name"`
	LastName      string `bun:"last_name" binding:"required" json:"last_name"`
	Role          string `bun:"role" binding:"required" json:"role"`
	IsEnabled     bool   `bun:"is_enabled,type:bool" json:"-" `
	IsDeleted     bool   `bun:"is_deleted,type:bool" json:"-"`
}

func AddUser(ctx context.Context, user *User) error {
	user.IsDeleted = false
	user.IsEnabled = true

	_, err := Db_GlobalVar.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		return fmt.Errorf("error adding user: %w", err)
	}

	return nil
}

func UserExists(ctx context.Context, username string) (bool, error) {
	var count int
	count, err := Db_GlobalVar.NewSelect().
		Model((*User)(nil)).
		Where("username = ?", username).
		Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetUserByUsername(ctx context.Context, username string) (*User, error) {
	user := new(User)

	err := Db_GlobalVar.NewSelect().Model(user).
		Where("username = ?", username).
		Where("is_deleted = ?", false).
		Scan(ctx)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("client with ClientID %s not found", username)
		}
		return nil, fmt.Errorf("error retrieving User cred with username %s: %w", username, err)
	}
	return user, nil
}

func GetAllUsers(ctx context.Context) ([]User, error) {
	var user []User
	err := Db_GlobalVar.NewSelect().Model(&user).Where("is_deleted = ?", false).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all Users: %w", err)
	}
	return user, nil
}

func UpdateUser(ctx context.Context, username string, updatedUser *User) (int64, error) {
	log.Debug().Msgf("Updating user with Username: %s\n", username)
	result, err := Db_GlobalVar.NewUpdate().
		Model(updatedUser).
		Where("is_deleted = ?", false).
		Where("username = ?", username).
		Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error updating User cred with username %s: %w", username, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error fetching rows affected: %w", err)
	}

	return rowsAffected, nil
}

func DeleteUser(ctx context.Context, username string) (int64, error) {
	log.Debug().Msgf("Deleting User with username: %s", username)

	result, err := Db_GlobalVar.NewUpdate().
		Model(&User{}).
		Set("is_deleted = ?", true).
		Where("username = ?", username).
		Exec(ctx)

	if err != nil {

		return 0, fmt.Errorf("error deleting client cred with ClientID %s: %w", username, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error fetching rows affected: %w", err)
	}

	return rowsAffected, nil
}
