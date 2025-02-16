package Ra

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func GetOneHandler(c echo.Context, db *pgxpool.Pool, ctx context.Context) error {
	resource := c.Param("resource")

	if err := CheckAllowedTables(resource); err != nil {
		return SendError(c, err)
	}

	id := c.Param("id")

	// FIXME: "ERROR: trailing junk after numeric literal at or near \"3f51fa8e\" (SQLSTATE 42601)"
	// This is an error converting uuid to json...
	dataQuery := fmt.Sprintf(`
		SELECT to_json(r) FROM (
			SELECT %s FROM %s where id = %s
		) r;`,
		strings.Join(
			AllowedTables[resource].ColumnsAllowed,
			", ",
		),
		resource,
		id,
	)

	row, err := db.Query(ctx, dataQuery)
	if err != nil {
		return SendError(c, err)
	}
	res := PgRowsToArrayMap(row)

	return c.JSON(http.StatusOK, res)
}
