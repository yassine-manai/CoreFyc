package migrations

import (
	"fyc/pkg/db"
	"reflect"
	"strings"

	"github.com/rs/zerolog/log"
)

type ColumnInfo struct {
	ColumnName string
	//DataType   string
}

func ValidateTableSchema(strct interface{}) error {
	db_validate := db.Db_GlobalVar.DB
	tableName := GetTableName(strct)

	query := `
		SELECT 
			column_name,
			data_type,
			is_nullable = 'YES' as is_nullable
		FROM information_schema.columns 
		WHERE table_name = $1
		ORDER BY ordinal_position;
	`

	log.Debug().Str("table", tableName).Msg("Executing query to fetch table schema")
	rows, err := db_validate.Query(query, tableName)
	if err != nil {
		log.Error().
			Err(err).
			Str("table", tableName).
			Msg("failed to query table schema")
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Error().Err(err).Str("table", tableName).Msg("failed to close rows")
		} else {
			log.Debug().Str("table", tableName).Msg("rows closed successfully")
		}
	}()

	if !rows.Next() {
		log.Warn().Str("table", tableName).Msg("Table not found in the database")
		return nil
	}

	dbColumns := make([]ColumnInfo, 0)

	var col ColumnInfo
	if err := rows.Scan(&col.ColumnName); err != nil {
		log.Error().
			Err(err).
			Str("table", tableName).
			Msg("failed to scan column info")
		return err
	}
	dbColumns = append(dbColumns, col)

	for rows.Next() {
		var col ColumnInfo
		err := rows.Scan(&col.ColumnName)
		if err != nil {
			log.Error().
				Err(err).
				Str("table", tableName).
				Msg("failed to scan column info")
			return err
		}
		dbColumns = append(dbColumns, col)
	}
	log.Debug().Int("Fetched columns length infos from database", len(dbColumns))
	if err := rows.Err(); err != nil {
		log.Error().Err(err).Str("table", tableName).Msg("error occurred while iterating over rows")
		return err
	}

	structType := reflect.TypeOf(strct)

	structColumns := make([]ColumnInfo, 0)
	log.Debug().Int("field_count", structType.NumField()).Msg("Processing struct fields for schema validation")

	for i := 1; i < structType.NumField(); i++ {
		field := structType.Field(i)

		tag := field.Tag.Get("bun")
		if tag == "" || tag == "-" {
			log.Debug().Str("field", field.Name).Msg("Skipping field without bun tag")
			continue
		}
		colName := strings.Split(tag, ",")[0]
		//dbType := goTypeToPostgres(field.Type)

		structColumns = append(structColumns, ColumnInfo{
			ColumnName: colName,
			//DataType:   dbType,
		})
	}

	log.Debug().Int("Processed struct columns for schema", len(structColumns))

	if len(dbColumns) != len(structColumns) {
		log.Warn().
			Str("table", tableName).
			Int("db_columns", len(dbColumns)).
			Int("struct_columns", len(structColumns)).
			Msg("column count mismatch")
		return err
	}

	for i, dbCol := range dbColumns {
		structCol := structColumns[i]

		log.Debug().
			Str("table", tableName).
			Int("position", i).
			Str("db_column", dbCol.ColumnName).
			Str("struct_column", structCol.ColumnName).
			Msg("Comparing columns")

		if dbCol.ColumnName != structCol.ColumnName {
			log.Warn().
				Str("table", tableName).
				Int("position", i).
				Str("db_column", dbCol.ColumnName).
				Str("struct_column", structCol.ColumnName).
				Msg("column name mismatch")
			return err
		}
	}

	log.Info().
		Str("table", tableName).
		Int("columns", len(dbColumns)).
		Msg("schema validation successful")

	return nil
}

/* // Helper function to check if types are compatible
func isCompatibleType(dbType, structType string) bool {
	compatibleTypes := map[string][]string{
		"integer": {"int", "integer"},
		"bigint":  {"bigint", "int64"},
		"text":    {"text", "varchar", "character varying"},
	}

	if compatible, ok := compatibleTypes[structType]; ok {
		for _, t := range compatible {
			if t == dbType {
				return true
			}
		}
	}

	isMatch := dbType == structType
	log.Debug().Str("db_type", dbType).Str("struct_type", structType).Bool("is_match", isMatch).Msg("Checking compatibility of types")
	return isMatch
} */

func GetTableName(structType interface{}) string {
	val := reflect.TypeOf(structType)
	log.Debug().Int("fields", val.NumField()).Msg("Retrieving table name from struct")

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		if bunTag := field.Tag.Get("bun"); bunTag != "" {
			for _, tagPart := range splitBunTags(bunTag) {
				if len(tagPart) > 6 && tagPart[:6] == "table:" {
					tableName := tagPart[6:]
					log.Debug().Str("struct_field", field.Name).Str("table_name", tableName).Msg("Found table name in struct tag")
					return tableName
				}
			}
		}
	}

	log.Warn().Msg("No table name found in struct tags")
	return ""
}

func splitBunTags(bunTag string) []string {
	return strings.Split(bunTag, ",")
}
