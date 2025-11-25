package exceptions

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"schneider.vip/problem"
)

func BuildErrAccountDocumentAlreadyExist(document string) *echo.HTTPError {
	return echo.NewHTTPError(
		http.StatusUnprocessableEntity,
		problem.Of(http.StatusUnprocessableEntity).Append(
			problem.Status(http.StatusUnprocessableEntity),
			problem.Detail(fmt.Sprintf("Account with document %s already exists", document)),
		),
	)
}

func BuildErrAccountNotFound(accountID string) *echo.HTTPError {
	return echo.NewHTTPError(
		http.StatusNotFound,
		problem.Of(http.StatusNotFound).Append(
			problem.Status(http.StatusNotFound),
			problem.Detail(fmt.Sprintf("account with uid %s not found", accountID)),
		),
	)
}
