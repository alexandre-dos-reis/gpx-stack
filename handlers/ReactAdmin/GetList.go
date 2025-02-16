package Ra

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// Ref: https://github.com/marmelab/react-admin/blob/master/packages/ra-data-simple-rest/README.md
type GetListRequest struct {
	Resource   string
	Sort       []string
	RangeStart int
	RangeEnd   int
	Filter     map[string]string
}

var AllowedTables map[string]struct{ ColumnsAllowed []string } = nil

func (req *GetListRequest) ToPgQuery() (*struct {
	DataQuery  string
	TotalQuery string
}, error,
) {
	start, end := req.RangeStart, req.RangeEnd
	limit := end - start + 1
	offset := start

	orderBy := ""
	if len(req.Sort) == 2 {
		column, direction := req.Sort[0], req.Sort[1]
		if direction != "ASC" && direction != "DESC" {
			return nil, errors.New("invalid sort direction")
		}
		orderBy = fmt.Sprintf("ORDER BY %s %s", column, direction)
	}

	// Build filtering query part
	var filterConditions []string
	for column, value := range req.Filter {
		filterConditions = append(filterConditions, fmt.Sprintf("%s = '%s'", column, value))
	}
	filterQuery := ""
	if len(filterConditions) > 0 {
		filterQuery = "WHERE " + strings.Join(filterConditions, " AND ")
	}

	// data dataQuery
	dataQuery := fmt.Sprintf(`
		SELECT row_to_json(r) FROM (
			SELECT %s FROM %s %s %s LIMIT %d OFFSET %d
		) r;`,
		strings.Join(
			AllowedTables[req.Resource].ColumnsAllowed,
			", ",
		),
		req.Resource,
		filterQuery,
		orderBy,
		limit,
		offset,
	)

	totalQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s %s;", req.Resource, filterQuery)

	return &struct {
		DataQuery  string
		TotalQuery string
	}{DataQuery: dataQuery, TotalQuery: totalQuery}, nil
}

func (req *GetListRequest) GetContentRange(dataLen int, total int) string {
	return fmt.Sprintf(
		"%s %d-%d/%d",
		req.Resource,
		req.RangeStart,
		req.RangeStart+dataLen-1,
		total,
	)
}

func (req *GetListRequest) Bind(
	c echo.Context,
) (*GetListRequest, error) {
	req.Resource = c.Param("resource")

	if err := CheckAllowedTables(req.Resource); err != nil {
		return nil, err
	}

	if err := BindValue(c.QueryParam("sort"), &req.Sort); err != nil {
		return nil, err
	}

	_range := []int{}
	if err := BindValue(c.QueryParam("range"), &_range); err != nil {
		return nil, err
	}
	if len(_range) < 2 {
		return nil, errors.New("Invalid range parameter")
	}

	req.RangeStart = _range[0]
	req.RangeEnd = _range[1]

	if err := BindValue(c.QueryParam("filter"), &req.Filter); err != nil {
		return nil, err
	}

	return req, nil
}

func GetListHandler(c echo.Context, db *pgxpool.Pool, ctx context.Context) error {
	request, err := (&GetListRequest{}).Bind(c)
	if err != nil {
		return SendError(c, err)
	}

	query, err := request.ToPgQuery()
	if err != nil {
		return SendError(c, err)
	}

	// begin transaction
	tx, err := db.Begin(ctx)
	if err != nil {
		return SendError(c, errors.New("Error with db query !"))
	}

	dataRows, err := tx.Query(ctx, query.DataQuery)
	if err != nil {
		tx.Rollback(ctx)
		return SendError(c, errors.New("Error with db query !"))
	}
	data := PgRowsToArrayMap(dataRows)
	dataRows.Close()

	// Count total query
	var total int
	totalErr := tx.QueryRow(ctx, query.TotalQuery).Scan(&total)
	if totalErr != nil {
		tx.Rollback(ctx)
		return SendError(c, errors.New("Error with db query !"))
	}

	// commit transaction
	if err = tx.Commit(ctx); err != nil {
		return SendError(c, errors.New("Error with db query !"))
	}

	c.Response().Header().Set("Content-Range", request.GetContentRange(len(data), total))
	c.Response().Header().Set("Access-Control-Expose-Headers", "Content-Range")
	return c.JSON(http.StatusOK, data)
}
