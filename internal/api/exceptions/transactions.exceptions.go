package exceptions

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"schneider.vip/problem"
)

func BuildErrTransactionInvalidAccount(accountID string) *echo.HTTPError {
	return echo.NewHTTPError(
		http.StatusUnprocessableEntity,
		problem.Of(http.StatusUnprocessableEntity).Append(
			problem.Status(http.StatusUnprocessableEntity),
			problem.Detail(fmt.Sprintf("account %s not found", accountID)),
		),
	)
}
