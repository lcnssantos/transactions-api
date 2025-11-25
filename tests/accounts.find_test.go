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

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type AccountFindTestSuite struct {
	suite.Suite
}

func (s *AccountFindTestSuite) TestShouldFindAccount() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	createResponse, err := testhttp.DoRequest[requests.CreateAccount, responses.Account](
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
	s.Require().Equal(http.StatusCreated, createResponse.Code)
	s.Require().NotNil(createResponse.Body)

	findResponse, err := testhttp.DoRequest[any, responses.Account](
		ctx,
		testhttp.Request[any]{
			Endpoint: fmt.Sprintf("http://localhost:%s/api/v1/accounts/%s", environmentVariables.Port, createResponse.Body.UID.String()),
			Method:   http.MethodGet,
		},
	)

	s.NoError(err)
	s.Equal(http.StatusOK, findResponse.Code)
	s.Require().NotNil(findResponse.Body)

	s.Equal(createResponse.Body.UID, findResponse.Body.UID)
	s.Equal(createResponse.Body.DocumentNumber, findResponse.Body.DocumentNumber)
}

func (s *AccountFindTestSuite) TestShouldReturnNotFoundWhenAccountDoesNotExist() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	randomID := uuid.New()

	findResponse, err := testhttp.DoRequest[any, map[string]any](
		ctx,
		testhttp.Request[any]{
			Method:   http.MethodGet,
			Endpoint: fmt.Sprintf("http://localhost:%s/api/v1/accounts/%s", environmentVariables.Port, randomID.String()),
		},
	)

	s.NoError(err)
	s.Equal(http.StatusNotFound, findResponse.Code)
}

func TestAccountFind(t *testing.T) {
	suite.Run(t, new(AccountFindTestSuite))
}
