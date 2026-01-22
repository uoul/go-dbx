package db

import (
	"database/sql"
	"reflect"
	"strings"
)

const (
	field_tag = "db"
)

func parseDbResult[T any](rows *sql.Rows) ([]T, error) {
	// Get column names from the result set
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	var result []T
	for rows.Next() {
		var item T
		// Create map of all fields from row
		fieldMap, err := createFieldMap(reflect.ValueOf(&item).Elem(), columns, "")
		if err != nil {
			return nil, err
		}
		// Create scan destinations using any typed interface
		scanDest := make([]any, len(columns))
		for i, col := range columns {
			if ptr, ok := fieldMap[col]; ok {
				scanDest[i] = ptr
			} else {
				// Skip unmapped fields into dummy variable
				var dummy any
				scanDest[i] = &dummy
			}
		}
		// Scan row
		if err := rows.Scan(scanDest...); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func createFieldMap(val reflect.Value, columns []string, prefix string) (map[string]any, error) {
	fieldMap := make(map[string]any)
	typ := val.Type()
	// Inspect all fields of type
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		fieldTag := fieldType.Tag.Get(field_tag)
		// Skip unexported fields
		if !field.CanSet() {
			continue
		}
		// Handle embedded structs
		if field.Kind() == reflect.Struct && fieldType.Anonymous {
			nestedMap, err := createFieldMap(field, columns, prefix)
			if err != nil {
				return nil, err
			}
			for k, v := range nestedMap {
				fieldMap[k] = v
			}
			continue
		}
		// Handle non-embedded nested structs
		if field.Kind() == reflect.Struct {
			nestedPrefix := fieldTag
			if nestedPrefix == "" {
				nestedPrefix = strings.ToLower(fieldType.Name)
			}
			// Add separator if there's already a prefix
			if prefix != "" {
				nestedPrefix = prefix + "_" + nestedPrefix
			}
			// Recursively process nested struct
			nestedMap, err := createFieldMap(field, columns, nestedPrefix)
			if err != nil {
				return nil, err
			}
			for k, v := range nestedMap {
				fieldMap[k] = v
			}
			continue
		}
		// Handle regular fields
		columnName := fieldTag
		if columnName == "" {
			columnName = strings.ToLower(fieldType.Name)
		}
		// Add prefix if exists
		if prefix != "" {
			columnName = prefix + "_" + columnName
		}
		// Add column to fieldmap
		fieldMap[columnName] = field.Addr().Interface()
	}
	return fieldMap, nil
}
