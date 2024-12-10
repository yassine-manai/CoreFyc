package db

import (
	"context"
	"fmt"
	"fyc/functions"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type ApiKey struct {
	bun.BaseModel `json:"-" bun:"table:api_key"`
	ID            int    `bun:"id,autoincrement,pk" json:"-"`
	ClientName    string `bun:"client_name" json:"client_name"`
	ClientID      string `bun:"client_id,unique" binding:"required" json:"client_id"`
	ClientSecret  string `bun:"client_secret,unique" binding:"required" json:"client_secret"`
	ApiKey        string `bun:"api_key" json:"-"`
	GrantType     string `bun:"grant_type" binding:"required" json:"grant_type"`
	FuzzyLogic    *bool  `bun:"fuzzy_logic,type:bool" json:"fuzzy_logic"`
	IsEnabled     bool   `bun:"is_enabled,type:bool" json:"-"`
	IsDeleted     bool   `bun:"is_deleted,type:bool" json:"-"`
	LastUpdated   string `bun:"last_update,type:timestamp" json:"-"`
}

type ApiKeyNoBind struct {
	bun.BaseModel `json:"-" bun:"table:api_key"`
	ID            int    `bun:"id,autoincrement" json:"_"`
	ClientName    string `bun:"client_name" json:"client_name"`
	ClientID      string `bun:"client_id,pk"  json:"client_id"`
	ClientSecret  string `bun:"client_secret,pk" json:"client_secret"`
	ApiKey        string `bun:"api_key" json:"-"`
	GrantType     string `bun:"grant_type" json:"grant_type"`
	FuzzyLogic    *bool  `bun:"fuzzy_logic,type:bool" json:"fuzzy_logic"`
	IsEnabled     *bool  `bun:"is_enabled,type:bool" json:"is_enabled"`
	IsDeleted     *bool  `bun:"is_deleted,type:bool" json:"-"`
	LastUpdated   string `bun:"last_update,type:timestamp" json:"-"`
}

type ApiKeyResponse struct {
	bun.BaseModel `json:"-" bun:"table:api_key"`
	ID            int     `bun:"id,autoincrement,pk" json:"-"`
	ClientName    string  `bun:"client_name" json:"client_name"`
	ClientID      string  `bun:"client_id,unique" json:"client_id"`
	ClientSecret  string  `bun:"client_secret,unique" json:"client_secret"`
	ApiKey        *string `bun:"api_key" json:"-"`
	GrantType     string  `bun:"grant_type" json:"grant_type"`
	FuzzyLogic    *bool   `bun:"fuzzy_logic,type:bool" json:"fuzzy_logic"`
	IsEnabled     bool    `bun:"is_enabled,type:bool" json:"is_enabled"`
	IsDeleted     bool    `bun:"is_deleted,type:bool" json:"-"`
	LastUpdated   string  `bun:"last_update,type:timestamp" json:"last_update"`
}

func GetAllDatas(ctx context.Context) (*ApiKeyResponse, error) {
	var api ApiKeyResponse
	err := Db_GlobalVar.NewSelect().Model(&api).Scan(ctx)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("error Getting all data")
		}
		return nil, fmt.Errorf("error Getting all dat  %w", err)
	}

	api.LastUpdated, _ = functions.ParseTimeData(api.LastUpdated)
	return &api, nil
}

func GetClientCredById(ctx context.Context, clientID string) (*ApiKeyResponse, error) {
	var api ApiKeyResponse

	err := Db_GlobalVar.NewSelect().Model(&api).
		Where("client_id = ?", clientID).
		Where("is_deleted = ?", false).
		Where("is_enabled = ?", true).
		Scan(ctx)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("client with ClientID %s not found", clientID)
		}
		return nil, fmt.Errorf("error retrieving client cred with ClientID %s: %w", clientID, err)
	}
	api.LastUpdated, _ = functions.ParseTimeData(api.LastUpdated)
	return &api, nil
}
func GetClientById(ctx context.Context, clientID string) (*ApiKeyResponse, error) {
	var api ApiKeyResponse

	err := Db_GlobalVar.NewSelect().Model(&api).
		Where("client_id = ?", clientID).
		Scan(ctx)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("client with ClientID %s not found", clientID)
		}
		return nil, fmt.Errorf("error retrieving client cred with ClientID %s: %w", clientID, err)
	}
	api.LastUpdated, _ = functions.ParseTimeData(api.LastUpdated)
	return &api, nil
}

func GetClientByStatus(ctx context.Context, status string) ([]ApiKeyResponse, error) {
	var client []ApiKeyResponse

	query := Db_GlobalVar.NewSelect().
		Model(&client).
		Where("is_deleted = ?", false)

	switch status {
	case "enabled":
		query.Where("is_enabled = ?", true)
	case "disabled":
		query.Where("is_enabled = ?", false)
	case "":
	default:
		return nil, fmt.Errorf("invalid status: %s, expected 'enabled', 'disabled', or 'all'", status)
	}

	// Execute the query
	if err := query.Scan(ctx); err != nil {
		return nil, fmt.Errorf("error getting client -- status: %s, err: %w", status, err)
	}

	// Parse `LastUpdated` for each zone
	for i := range client {
		client[i].LastUpdated, _ = functions.ParseTimeData(client[i].LastUpdated)
	}

	return client, nil
}

func GetClientCredByIdFalse(ctx context.Context, clientID string) (*ApiKeyResponse, error) {
	var api ApiKeyResponse

	err := Db_GlobalVar.NewSelect().Model(&api).
		Where("client_id = ?", clientID).
		Where("is_deleted = ?", false).
		//Where("is_enabled = ?", true).
		Scan(ctx)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("client with ClientID %s not found", clientID)
		}
		return nil, fmt.Errorf("error retrieving client cred with ClientID %s: %w", clientID, err)
	}
	api.LastUpdated, _ = functions.ParseTimeData(api.LastUpdated)
	return &api, nil
}
func GetClientCredBySecret(ctx context.Context, clientSecret string) (*ApiKeyResponse, error) {
	var api ApiKeyResponse

	err := Db_GlobalVar.NewSelect().
		Model(&api).
		Where("client_secret = ?", clientSecret).
		Where("is_deleted = ?", false).
		Where("is_enabled = ?", true).
		Scan(ctx)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("client with Client Secret %s not found", clientSecret)
		}
		return nil, fmt.Errorf("error retrieving client cred with Client Secret %s: %w", clientSecret, err)
	}
	api.LastUpdated, _ = functions.ParseTimeData(api.LastUpdated)
	return &api, nil
}

func GetAllClientCred(ctx context.Context) ([]ApiKeyResponse, error) {
	var apm []ApiKeyResponse
	err := Db_GlobalVar.NewSelect().
		Model(&apm).
		Where("is_deleted = ?", false).
		//Where("is_enabled = ?", true).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all Client Credentials: %w", err)
	}

	for i := range apm {
		apm[i].LastUpdated, _ = functions.ParseTimeData(apm[i].LastUpdated)
	}
	return apm, nil
}

func GetAllClientDatas(ctx context.Context) ([]ApiKey, error) {
	var api []ApiKey
	err := Db_GlobalVar.NewSelect().Model(&api).Scan(ctx)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("error Getting all data")
		}
		return nil, fmt.Errorf("error Getting all dat  %w", err)
	}

	for i := range api {
		api[i].LastUpdated, _ = functions.ParseTimeData(api[i].LastUpdated)
	}
	return api, nil
}

func AddClientCred(ctx context.Context, apimgnt *ApiKey) error {
	log.Debug().Str("Client API Added AT:", functions.GetFormatedLocalTime())
	apimgnt.LastUpdated = functions.GetFormatedLocalTime()
	apimgnt.IsDeleted = false
	apimgnt.IsEnabled = true
	//apimgnt.GrantType = "client_credentials"

	_, err := Db_GlobalVar.NewInsert().Model(apimgnt).Exec(ctx)
	if err != nil {
		return fmt.Errorf("error adding api_cred: %w", err)
	}

	LoadClientDataList()
	LoadClientsApi()
	LoadClientlist()
	return nil
}

func UpdateClientCred(ctx context.Context, clientID string, updatedClientCred *ApiKeyNoBind) (int64, error) {
	log.Debug().Str("Client Added AT:", functions.GetFormatedLocalTime()).Str("Client ID:", clientID)
	updatedClientCred.LastUpdated = functions.GetFormatedLocalTime()

	result, err := Db_GlobalVar.NewUpdate().
		Model(updatedClientCred).
		Where("client_id = ?", clientID).
		Where("is_deleted = ?", false).
		OmitZero().
		Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error updating client cred with ClientID %s: %w", clientID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error fetching rows affected: %w", err)
	}

	LoadClientDataList()
	LoadClientsApi()
	LoadClientlist()
	return rowsAffected, nil
}

func DeleteClientCred(ctx context.Context, clientID string) (int64, error) {
	log.Debug().Str("At ", functions.GetFormatedLocalTime()).Msgf("Deleting Client with ClientID: %s", clientID)

	result, err := Db_GlobalVar.NewUpdate().
		Model(&ApiKeyResponse{}).
		Where("client_id = ?", clientID).
		Set("is_deleted = ?", true).
		Set("is_enabled = ?", false).
		Set("last_update = ?", functions.GetFormatedLocalTime()).
		Exec(ctx)

	if err != nil {

		return 0, fmt.Errorf("error deleting client cred with ClientID %s: %w", clientID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error fetching rows affected: %w", err)
	}
	LoadClientDataList()
	LoadClientsApi()
	LoadClientlist()
	return rowsAffected, nil
}

func StoreToken(ctx context.Context, clientID string, token string) (int64, error) {
	res, err := Db_GlobalVar.NewUpdate().
		Model(&ApiKeyResponse{}).
		Where("is_deleted = ?", false).
		Where("is_enabled = ?", true).
		Where("client_id = ?", clientID).
		Set("api_key = ?", token).
		Exec(ctx)

	if err != nil {
		return 0, fmt.Errorf("error updating token with clientID %s: %w", clientID, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Info().Str("Client ID", clientID).Int("Rows Affected ", int(rowsAffected)).Msg("Changed token")

	return rowsAffected, nil
}

/* func GetTokenByClientID(ctx context.Context, clientID string) (*ApiKeyResponse, error) {
	api := new(ApiKeyResponse)

	err := Db_GlobalVar.NewSelect().Model(api).
		Where("client_id = ?", clientID).
		Where("is_deleted = ?", false).
		Where("is_enabled = ?", true).
		Scan(ctx)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("client with ID %s not found", clientID)
		}
		return nil, fmt.Errorf("error retrieving api key  with ClientID %s: %w", clientID, err)
	}
	return api, nil
}
*/
