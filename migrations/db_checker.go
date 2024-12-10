package migrations

import (
	"reflect"
	"strings"

	"github.com/rs/zerolog/log"
)

func GetStructColomns(strct interface{}) {
	structType := reflect.TypeOf(strct)

	log.Debug().Interface("///// STRCT ///// ", strct).Send()

	structColumns := make([]ColumnInfo, 0)
	log.Debug().Int("field_count", structType.NumField()).Msg("Processing struct fields for schema validation")

	for i := 0; i < structType.NumField(); i++ {
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

	log.Debug().Interface("Processed struct columns for schema", structColumns).Send()

}

func GetDBColomns(tableName string) {
	//db_validate string,
	log.Debug().Str("Name", tableName).Msg("Get DB COLS")

}

/* func goTypeToPostgres(t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.Int, reflect.Int32:
		return "integer"
	case reflect.Int64:
		return "bigint"
	case reflect.Float32:
		return "real"
	case reflect.Float64:
		return "double precision"
	case reflect.String:
		return "text"
	case reflect.Bool:
		return "boolean"
	case reflect.Interface:
		return "jsonb"
	}

	log.Warn().Str("type", t.String()).Msg("unknown type, defaulting to 'unknown'")
	return "unknown"
}
*/
