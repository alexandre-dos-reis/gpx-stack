package Ra

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func CheckAllowedTables(resource string) error {
	if _, exists := AllowedTables[resource]; !exists {
		return errors.New("The resource doesn't exists !")
	}
	return nil
}

func BindValue[T any](param string, value *T) error {
	if param != "" {
		if err := json.Unmarshal([]byte(param), &value); err != nil {
			return err
		}
	}
	return nil
}

func SendError(c echo.Context, err error) error {
	// TODO: See ra manual for error formatting
	return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
}

func PgRowsToArrayMap(rows pgx.Rows) []map[string]interface{} {
	tableData := make([]map[string]interface{}, 0)

	for rows.Next() {
		var jsonData []byte
		if err := rows.Scan(&jsonData); err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}

		// Unmarshal JSON into a map
		var entry map[string]interface{}
		if err := json.Unmarshal(jsonData, &entry); err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			continue
		}

		tableData = append(tableData, entry)
	}

	return tableData
}

func PgRowToMap(row pgx.Row) (map[string]interface{}, error) {
	var jsonData string // PostgreSQL returns a JSON string

	// Scan directly into a string
	if err := row.Scan(&jsonData); err != nil {
		return nil, fmt.Errorf("error scanning row: %w", err)
	}

	// Unmarshal JSON string into a map
	var entry map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &entry); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	return entry, nil
}
