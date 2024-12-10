package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type ErrorMessage struct {
	bun.BaseModel `json:"-" bun:"table:errors"`
	Code          int               `bun:"code" json:"code"`
	Messages      map[string]string `bun:"messages,type:jsonb" json:"messages" swaggertype:"object"`
}

func CreateErrorMessage(ctx context.Context, errMsg *ErrorMessage) error {
	_, err := Db_GlobalVar.NewInsert().Model(errMsg).Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to insert error message into database")
		return err
	}
	log.Info().Int("code", errMsg.Code).Msg("Successfully inserted error message")
	return nil
}

// GetErrorMessageByFilter fetches an error message by code and language from the database
func GetErrorMessageByFilter(ctx context.Context, code int, language string) (ErrorMessage, error) {
	var errMsg ErrorMessage
	err := Db_GlobalVar.NewSelect().
		Model(&errMsg).
		Where("code = ?", code).
		Where("messages ->> ? IS NOT NULL", language).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn().Int("code", code).Str("language", language).Msg("No error message found for the given code and language")
		} else {
			log.Error().Err(err).Int("code", code).Str("language", language).Msg("Failed to fetch error message from database")
		}
		return errMsg, err
	}

	log.Info().Int("code", code).Str("language", language).Msg("Successfully fetched error message")
	return errMsg, nil
}

func GetErrorMessageByCode(ctx context.Context, code int) (ErrorMessage, error) {
	var errMsg ErrorMessage
	err := Db_GlobalVar.NewSelect().
		Model(&errMsg).
		Where("code = ?", code).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn().Int("code", code).Msg("No error message found for the given code")
		} else {
			log.Error().Err(err).Int("code", code).Msg("Failed to fetch error message from database")
		}
		return errMsg, err
	}

	log.Info().Int("code", code).Msg("Successfully fetched error message")
	return errMsg, nil
}

func GetErrorMessage(ctx context.Context) ([]ErrorMessage, error) {
	var errMsg []ErrorMessage
	err := Db_GlobalVar.NewSelect().
		Model(&errMsg).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn().Msg("No error messages found")
		} else {
			log.Error().Err(err).Msg("Failed to fetch error messages from database")
		}
		return errMsg, err
	}

	log.Info().Msg("Successfully fetched error messages")
	return errMsg, nil
}

// UpdateErrorMessage updates an error message in the database
func UpdateErrorMessage(ctx context.Context, errMsg *ErrorMessage) error {
	_, err := Db_GlobalVar.NewUpdate().
		Model(errMsg).
		Set("messages = ?", errMsg.Messages).
		Where("code = ?", errMsg.Code).
		Exec(ctx)

	if err != nil {
		log.Error().Err(err).Int("code", errMsg.Code).Msg("Failed to update error message in database")
		return err
	}

	log.Info().Int("code", errMsg.Code).Msg("Successfully updated error message")
	return nil
}

// DeleteErrorMessage removes a specific language from the messages field of an error message by code
func DeleteErrorMessage(ctx context.Context, code int, language string) (int64, error) {
	res, err := Db_GlobalVar.NewUpdate().
		Model((*ErrorMessage)(nil)).
		Set("messages = messages - ?", language).
		Where("code = ?", code).
		Exec(ctx)

	if err != nil {
		log.Error().Err(err).Int("code", code).Str("language", language).Msg("Failed to delete language from error message")
		return 0, err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		log.Warn().Int("code", code).Str("language", language).Msg("No error message found or no language removed")
	} else {
		log.Info().Int("code", code).Str("language", language).Int64("rowsAffected", rowsAffected).Msg("Successfully deleted language from error message")
	}

	return rowsAffected, nil
}
