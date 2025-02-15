package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var allowedTables = map[string]struct{ columnsAllowed []string }{
	"products": {columnsAllowed: []string{"*"}},
}

// Ref: https://github.com/marmelab/react-admin/blob/master/packages/ra-data-simple-rest/README.md
type GetListRequest struct {
	Resource   string
	Sort       []string
	RangeStart int
	RangeEnd   int
	Filter     map[string]string
}
type GetListResponse []any

func BindValue[T any](param string, value *T) error {
	if param != "" {
		if err := json.Unmarshal([]byte(param), &value); err != nil {
			return err
		}
	}
	return nil
}

func (req *GetListRequest) toPgQuery() (*struct {
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
			allowedTables[req.Resource].columnsAllowed,
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

func (req *GetListRequest) bind(c echo.Context) (*GetListRequest, error) {
	req.Resource = c.Param("resource")

	if _, exists := allowedTables[req.Resource]; !exists {
		return nil, errors.New("The resource doesn't exists !")
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

func JsonBadRequest(c echo.Context, err error) error {
	// TODO: See ra manual for error formatting
	return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
}

func (h *Handlers) adminHandlers() {
	h.echo.Use(middleware.Logger())
	// useless if admin is served by the same server...
	h.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPut,
			http.MethodPatch,
			http.MethodPost,
			http.MethodDelete,
		},
	}))
	h.echo.Debug = true

	g := h.echo.Group("/admin")

	g.GET("/:resource", func(c echo.Context) error {
		request, err := (&GetListRequest{}).bind(c)
		if err != nil {
			return JsonBadRequest(c, err)
		}

		query, err := request.toPgQuery()
		if err != nil {
			return JsonBadRequest(c, err)
		}

		// begin transaction
		tx, err := h.db.Begin(h.ctx)
		if err != nil {
			return JsonBadRequest(c, errors.New("Error with db query !"))
		}

		dataRows, err := tx.Query(h.ctx, query.DataQuery)
		if err != nil {
			tx.Rollback(h.ctx)
			return JsonBadRequest(c, errors.New("Error with db query !"))
		}
		data := pgRowsToArrayMap(dataRows)
		dataRows.Close()

		// Count total query
		var total int
		totalErr := tx.QueryRow(h.ctx, query.TotalQuery).Scan(&total)
		if totalErr != nil {
			tx.Rollback(h.ctx)
			return JsonBadRequest(c, errors.New("Error with db query !"))
		}

		// commit transaction
		if err = tx.Commit(h.ctx); err != nil {
			return JsonBadRequest(c, errors.New("Error with db query !"))
		}

		c.Response().Header().Set("Content-Range", request.GetContentRange(len(data), total))
		c.Response().Header().Set("Access-Control-Expose-Headers", "Content-Range")
		return c.JSON(http.StatusOK, data)
	})
}

func pgRowsToArrayMap(rows pgx.Rows) []map[string]interface{} {
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
