package db

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

type UserAudit struct {
	bun.BaseModel `json:"-" bun:"table:user_audit"`
	ID            int                    `bun:"id,autoincrement" json:"id"`
	UserID        int                    `bun:"user_id,pk" binding:"required" json:"user_id"`
	ActionDate    string                 `bun:"action_date,pk,type:timestamp" binding:"required" json:"action_date"`
	OldValue      map[string]interface{} `bun:"old_value,type:jsonb" binding:"required" json:"old_value" swaggertype:"object"`
	NewValue      map[string]interface{} `bun:"new_value,type:jsonb" binding:"required" json:"new_value" swaggertype:"object"`
	Module        string                 `bun:"module" binding:"required" json:"module"`
}

func CreateUserAudit(ctx context.Context, userAudit *UserAudit) error {
	_, err := Db_GlobalVar.NewInsert().Model(userAudit).Exec(ctx)
	if err != nil {
		return fmt.Errorf("error adding userAudit: %w", err)
	}
	return nil
}

func GetUserAuditById(ctx context.Context, UserAuditID int) (*UserAudit, error) {
	var UserAudit UserAudit

	err := Db_GlobalVar.NewSelect().Model(&UserAudit).
		Where("user_id = ?", UserAuditID).
		Scan(ctx)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("UserAudit with ID %d not found", UserAuditID)
		}
		return nil, fmt.Errorf("error retrieving UserAudit with UserAuditID %d: %w", UserAuditID, err)
	}
	return &UserAudit, nil
}

func GetAllUserAudits(ctx context.Context) ([]UserAudit, error) {
	var UserAudits []UserAudit
	err := Db_GlobalVar.NewSelect().Model(&UserAudits).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all UserAudits: %w", err)
	}
	return UserAudits, nil
}

func UpdateUserAudit(ctx context.Context, UserAuditID int, updatedUserAudit *UserAudit) (int64, error) {
	//log.Debug().Str("Time", global.FormattedTime).Msgf("Updating UserAudit with ID: %d\n", UserAuditID)
	//updatedUserAudit.ActionDate = global.FormattedTime

	result, err := Db_GlobalVar.NewUpdate().
		Model(updatedUserAudit).
		Where("user_id = ?", UserAuditID).
		Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error updating UserAudit with UserAuditID %d: %w", UserAuditID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error fetching rows affected: %w", err)
	}

	return rowsAffected, nil
}

func DeleteUserAudit(ctx context.Context, UserAuditID int) (int64, error) {
	//log.Debug().Str("Time ", global.FormattedTime).Msgf("Deleting UserAudit with UserAuditID: %d", UserAuditID)

	result, err := Db_GlobalVar.NewDelete().
		Model(&UserAudit{}).
		Where("user_id = ?", UserAuditID).
		Exec(ctx)

	if err != nil {
		return 0, fmt.Errorf("error deleting UserAudit with UserAuditID %d: %w", UserAuditID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error fetching rows affected: %w", err)
	}

	return rowsAffected, nil
}
