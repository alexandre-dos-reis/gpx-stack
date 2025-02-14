package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type GetListRequest struct {
	Resource string
	Sort     []string
	Range    []int
	Filter   map[string]string
}

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
	h.echo.Use(middleware.CORS())
	h.echo.Debug = true

	h.echo.Use(middleware.Logger())

	g := h.echo.Group("/admin")

	g.GET("/:resource", func(c echo.Context) error {
		request, err := (&GetListRequest{}).bind(c)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, "bad request")
		}
		return c.JSON(http.StatusOK, &request)
	})
}
