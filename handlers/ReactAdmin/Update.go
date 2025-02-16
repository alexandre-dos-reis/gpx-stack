package Ra

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func UpdateHandler(c echo.Context, db *pgxpool.Pool, ctx context.Context) error {
	resource := c.Param("resource")

	if err := CheckAllowedTables(resource); err != nil {
		return SendError(c, err)
	}

	id := c.Param("id")

	var data map[string]any
	if err := c.Bind(&data); err != nil {
		return SendError(c, err)
	}

	// FIXME: Resource seemd to be present in the columns, rm it !

	columns := []string{}
	values := []string{}

	for k, v := range data {
		columns = append(columns, k)
		values = append(values, fmt.Sprintf("%v", v))
	}

	insertedRow := db.QueryRow(ctx, `
		SELECT to_json(r) FROM (
			UPDATE %s SET (%s) = (%s) WHERE id = '%s' RETURNING *
		) r;`,
		resource,
		strings.Join(columns, ", "),
		strings.Join(values, ", "),
		id,
	)

	res, err := PgRowToMap(insertedRow)
	if err != nil {
		return SendError(c, err)
	}

	return c.JSON(http.StatusOK, res)
}
