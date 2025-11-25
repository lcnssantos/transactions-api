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

type Account struct {
	accountService services.Account
}

func NewAccount(
	accountService services.Account,
) Account {
	return Account{
		accountService: accountService,
	}
}

func (a Account) createAccount(c echo.Context) error {
	body, err := rest.Bind[requests.CreateAccount](c)

	if err != nil {
		return err
	}

	account, err := a.accountService.Create(c.Request().Context(), body.Domain())

	if err != nil {
		if errors.Is(err, repositories.ErrAccountDocumentAlreadyExist) {
			return exceptions.BuildErrAccountDocumentAlreadyExist(body.DocumentNumber)

		}
		return err
	}

	return c.JSON(http.StatusCreated, responses.Account{}.FromDomain(account))
}

func (a Account) findAccount(c echo.Context) error {
	body, err := rest.Bind[requests.FindAccount](c)

	if err != nil {
		return err
	}

	account, err := a.accountService.FindByUUID(c.Request().Context(), body.AccountID)

	if err != nil {
		if errors.Is(err, repositories.ErrAccountNotFound) {
			return exceptions.BuildErrAccountNotFound(body.AccountID.String())
		}

		return err
	}

	return c.JSON(http.StatusOK, responses.Account{}.FromDomain(account))
}

func (a Account) RegisterRoutes(e *echo.Group) {
	e.POST("", a.createAccount)
	e.GET("/:account_id", a.findAccount)
}
