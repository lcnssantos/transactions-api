package main

import (
	"context"
	"fmt"
	"os"
	"time"
	"transactions/internal/api"
	"transactions/internal/env"
	"transactions/internal/lib/environment"
	"transactions/internal/lib/logging"
	"transactions/internal/lib/postgres"
	"transactions/internal/lib/validator"
	"transactions/internal/repositories"
	"transactions/internal/services"

	"github.com/labstack/echo/v4"

	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func main() {
	ctx := context.TODO()

	logging.Init(os.Getenv("LOG_LEVEL"))

	traceExporter, err := otlptracehttp.New(
		ctx,
	)

	if err != nil {
		logging.Panic(ctx, err).Msg("error to start otel traceExporter")
		return
	}

	resource, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("transactions.api"),
			semconv.ServiceVersionKey.String("v1.0.0"),
		),
	)

	if err != nil {
		logging.Panic(ctx, err).Msg("error to start otel resource")
		return
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(resource),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	environmentVariables, err := environment.LoadEnvironment[env.Environment]()

	if err != nil {
		logging.Panic(ctx, err).Msg("error to load environment variables")
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

	migrationsPath := fmt.Sprintf("file://%s", "./migrations")

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

	err = app.Start(fmt.Sprintf(":%s", environmentVariables.Port))

	if err != nil {
		logging.Panic(ctx, err).Msg("error to start http server")
		return
	}
}
