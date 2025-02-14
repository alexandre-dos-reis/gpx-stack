package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/a-h/templ-examples/hello-world/database/repository"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type GetListRequest struct {
	Resource string
	Sort     []string
	Range    []int
	Filter   map[string]string
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

func (req *GetListRequest) bind(c echo.Context) (*GetListRequest, error) {
	req.Resource = c.Param("resource")

	if err := BindValue(c.QueryParam("sort"), &req.Sort); err != nil {
		return nil, err
	}
	if err := BindValue(c.QueryParam("range"), &req.Range); err != nil {
		return nil, err
	}
	if err := BindValue(c.QueryParam("filter"), &req.Filter); err != nil {
		return nil, err
	}

	return req, nil
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
			return c.JSON(http.StatusBadRequest, "bad request")
		}

		resourceTableMapping := map[string]struct {
			table      string
			structType interface{}
		}{
			"products": {table: "products", structType: repository.Product{}},
		}

		table := resourceTableMapping[request.Resource]

		contentRange := fmt.Sprintf("%s %d-%d/%d", request.Resource, 0, 0, 0)
		c.Response().Header().Set("Content-Range", contentRange)
		c.Response().Header().Set("Access-Control-Expose-Headers", "Content-Range")
		return c.JSON(http.StatusOK, &GetListResponse{})
	})
}
