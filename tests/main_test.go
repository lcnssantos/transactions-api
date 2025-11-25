package tests

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
	"transactions/internal/api"
	"transactions/internal/env"
	"transactions/internal/lib/environment"
	"transactions/internal/lib/logging"
	"transactions/internal/lib/postgres"
	"transactions/internal/lib/validator"
	"transactions/internal/repositories"
	"transactions/internal/services"
	"transactions/tests/internal/utils"

	"github.com/labstack/echo/v4"
)

var environmentVariables env.Environment

func TestMain(m *testing.M) {
	ctx := context.TODO()

	logging.Init(os.Getenv("LOG_LEVEL"))

	var err error

	environmentVariables, err = environment.LoadEnvironment[env.Environment]("../.env")

	if err != nil {
		logging.Panic(ctx, err).Msg("error to load environment variables")
		return
	}

	err = utils.AssertTestDatabase(
		environmentVariables.DatabaseHost,
		environmentVariables.DatabaseUser,
		environmentVariables.DatabasePass,
		environmentVariables.DatabaseName,
		environmentVariables.DatabasePort,
		environmentVariables.DatabaseSSLMode,
	)

	environmentVariables.DatabaseName = fmt.Sprintf("%s_test", environmentVariables.DatabaseName)

	if err != nil {
		logging.Panic(ctx, err).Msg("error to assert database")
		return
	}

	poolConfig := postgres.NewPoolConfig(environmentVariables.DatabasePoolMinimum, environmentVariables.DatabasePoolMaximum, time.Second)

	pg := postgres.New(
		postgres.NewConfig(
			environmentVariables.DatabaseHost,
			environmentVariables.DatabasePort,
			environmentVariables.DatabaseUser,
			environmentVariables.DatabasePass,
			environmentVariables.DatabaseName,
			environmentVariables.DatabaseSSLMode,
		),
	).WithPoolConfig(poolConfig)

	err = pg.Connect()

	if err != nil {
		logging.Panic(ctx, err).Msg("error to connect to database")
		return
	}

	migrationsPath := fmt.Sprintf("file://%s", "../migrations")

	err = pg.MigrateDown(migrationsPath)

	if err != nil {
		logging.Panic(ctx, err).Msg("error to migrate down database")
		return
	}

	err = pg.MigrateUp(migrationsPath)

	if err != nil {
		logging.Panic(ctx, err).Msg("error to migrate database")
		return
	}

	app := echo.New()

	app.Validator, err = validator.NewCustomValidator(
		validator.Locale_en,
	)

	if err != nil {
		logging.Panic(ctx, err).Msg("error to create validator")
	}

	accountRepository := repositories.NewAccount(pg.DB())
	transactionRepository := repositories.NewTransaction(pg.DB())

	accountService := services.NewAccount(accountRepository)
	transactionService := services.NewTransaction(transactionRepository)

	api.RoutesFactory(accountService, transactionService)(app.Group("/api"))

	go func() {
		err = app.Start(fmt.Sprintf(":%s", environmentVariables.Port))

		if err != nil {
			logging.Panic(ctx, err).Msg("error to start http server")
			return
		}
	}()

	m.Run()
}
