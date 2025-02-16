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

	dataQuery := fmt.Sprintf(`
		SELECT to_json(r) FROM (
			SELECT %s FROM %s where id = '%s'
		) r;`,
		strings.Join(
			AllowedTables[resource].ColumnsAllowed,
			", ",
		),
		resource,
		id,
	)

	row := db.QueryRow(ctx, dataQuery)
	res, err := PgRowToMap(row)
	if err != nil {
		return SendError(c, err)
	}

	return c.JSON(http.StatusOK, res)
}
