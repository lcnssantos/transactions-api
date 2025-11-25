package controllers

import (
	"errors"
	"net/http"
	"transactions/internal/api/exceptions"
	"transactions/internal/api/requests"
	"transactions/internal/api/responses"
	"transactions/internal/lib/rest"
	"transactions/internal/repositories"
	"transactions/internal/services"

	"github.com/labstack/echo/v4"
)

type Transaction struct {
	transactionService services.Transaction
	accountService     services.Account
}

func NewTransaction(transactionService services.Transaction, accountService services.Account) Transaction {
	return Transaction{
		transactionService: transactionService,
		accountService:     accountService,
	}
}

func (t Transaction) createTransaction(c echo.Context) error {
	ctx := c.Request().Context()

	body, err := rest.Bind[requests.CreateTransaction](c)

	if err != nil {
		return err
	}

	account, err := t.accountService.FindByUUID(ctx, body.AccountID)

	if err != nil {
		if errors.Is(err, repositories.ErrAccountNotFound) {
			return exceptions.BuildErrTransactionInvalidAccount(body.AccountID.String())
		}

		return err
	}

	transaction := body.Domain(account)

	transaction, err = t.transactionService.Save(ctx, transaction)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, responses.Transaction{}.FromDomain(transaction))
}

func (t Transaction) RegisterRoutes(e *echo.Group) {
	e.POST("", t.createTransaction)
}
