package tests

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
	"transactions/internal/api/requests"
	"transactions/internal/api/responses"
	"transactions/tests/internal/testhttp"
	"transactions/tests/internal/utils"

	"github.com/stretchr/testify/suite"
)

type AccountCreationTestSuite struct {
	suite.Suite
}

func (s *AccountCreationTestSuite) TestShouldCreateAccount() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	response, err := testhttp.DoRequest[requests.CreateAccount, responses.Account](
		ctx,
		testhttp.Request[requests.CreateAccount]{
			Method:   http.MethodPost,
			Endpoint: fmt.Sprintf("http://localhost:%s/api/v1/accounts", environmentVariables.Port),
			Body: &requests.CreateAccount{
				DocumentNumber: utils.GenerateDocument(11),
			},
		},
	)

	s.NoError(err)

	s.Equal(http.StatusCreated, response.Code)
	s.Equal(response.Body.DocumentNumber, response.Body.DocumentNumber)
}

func (s *AccountCreationTestSuite) TestShouldReturnErrorWhenDocumentAlreadyExists() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	document := utils.GenerateDocument(11)

	_, err := testhttp.DoRequest[requests.CreateAccount, responses.Account](
		ctx,
		testhttp.Request[requests.CreateAccount]{
			Method:   http.MethodPost,
			Endpoint: fmt.Sprintf("http://localhost:%s/api/v1/accounts", environmentVariables.Port),
			Body: &requests.CreateAccount{
				DocumentNumber: document,
			},
		},
	)

	s.NoError(err)

	response, err := testhttp.DoRequest[requests.CreateAccount, responses.Account](
		ctx,
		testhttp.Request[requests.CreateAccount]{
			Method:   http.MethodPost,
			Endpoint: fmt.Sprintf("http://localhost:%s/api/v1/accounts", environmentVariables.Port),
			Body: &requests.CreateAccount{
				DocumentNumber: document,
			},
		},
	)

	s.NoError(err)

	s.Equal(http.StatusUnprocessableEntity, response.Code)
}

func TestAccountCreation(t *testing.T) {
	suite.Run(t, new(AccountCreationTestSuite))
}
