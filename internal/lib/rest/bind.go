package rest

import (
	"github.com/labstack/echo/v4"
)

func Bind[T any](c echo.Context) (T, error) {
	var t T

	if err := c.Bind(&t); err != nil {
		return t, err
	}

	if err := c.Validate(&t); err != nil {
		return t, err
	}

	return t, nil
}
