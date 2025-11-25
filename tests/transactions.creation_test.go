package tests

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
	"transactions/internal/api/requests"
	"transactions/internal/api/responses"
	"transactions/internal/domain"
	"transactions/tests/internal/testhttp"
	"transactions/tests/internal/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type TransactionCreationTestSuite struct {
	suite.Suite
}

func (s *TransactionCreationTestSuite) createAccount(ctx context.Context) responses.Account {
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

	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, response.Code)
	s.Require().NotNil(response.Body)

	return *response.Body
}

func (s *TransactionCreationTestSuite) TestShouldCreateTransaction() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	account := s.createAccount(ctx)
	externalID := uuid.New()

	response, err := testhttp.DoRequest[requests.CreateTransaction, responses.Transaction](
		ctx,
		testhttp.Request[requests.CreateTransaction]{
			Method:   http.MethodPost,
			Endpoint: fmt.Sprintf("http://localhost:%s/api/v1/transactions", environmentVariables.Port),
			Body: &requests.CreateTransaction{
				ExternalID:    externalID,
				AccountID:     account.UID,
				OperationType: domain.TransactionOperationTypePurchase,
				Amount:        1000,
			},
		},
	)

	s.NoError(err)
	s.Equal(http.StatusOK, response.Code)
	s.Require().NotNil(response.Body)

	s.Equal(domain.TransactionOperationTypePurchase, response.Body.OperationType)
	s.Equal(int64(1000), response.Body.Amount)
	s.Equal(externalID, response.Body.ExternalID)
	s.Equal(account.UID, response.Body.Account.UID)
	s.Equal(account.DocumentNumber, response.Body.Account.DocumentNumber)
	s.False(response.Body.EventDate.IsZero())
}

func (s *TransactionCreationTestSuite) TestShouldReturnErrorWhenAccountDoesNotExist() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	randomAccount := uuid.New()

	response, err := testhttp.DoRequest[requests.CreateTransaction, map[string]any](
		ctx,
		testhttp.Request[requests.CreateTransaction]{
			Method:   http.MethodPost,
			Endpoint: fmt.Sprintf("http://localhost:%s/api/v1/transactions", environmentVariables.Port),
			Body: &requests.CreateTransaction{
				ExternalID:    uuid.New(),
				AccountID:     randomAccount,
				OperationType: domain.TransactionOperationTypePurchase,
				Amount:        500,
			},
		},
	)

	s.NoError(err)
	s.Equal(http.StatusUnprocessableEntity, response.Code)
}

func TestTransactionCreation(t *testing.T) {
	suite.Run(t, new(TransactionCreationTestSuite))
}
