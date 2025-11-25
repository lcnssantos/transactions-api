package api

import (
	"transactions/internal/api/controllers"
	"transactions/internal/api/middlewares"
	"transactions/internal/services"

	"github.com/labstack/echo/v4"
)

func RoutesFactory(
	accountService services.Account,
	transactionService services.Transaction,
) func(f *echo.Group) {
	accountsController := controllers.NewAccount(accountService)
	transactionsController := controllers.NewTransaction(transactionService, accountService)

	return func(f *echo.Group) {
		f.Use(middlewares.Observability())

		accountsController.RegisterRoutes(f.Group("/v1/accounts"))
		transactionsController.RegisterRoutes(f.Group("/v1/transactions"))
	}
}
