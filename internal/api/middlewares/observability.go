package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"transactions/internal/lib/logging"
	"transactions/internal/lib/rest"

	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type responseBodyCapturer struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (w *responseBodyCapturer) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Observability() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request := c.Request()

			path := c.Request().URL.Path
			route := fmt.Sprintf("%s %s", request.Method, path)

			otelCtx := otel.GetTextMapPropagator().
				Extract(request.Context(), propagation.HeaderCarrier(c.Request().Header))

			requestTime := time.Now()

			requestBodyBytes, _ := io.ReadAll(request.Body)
			request.Body.Close()

			request.Body = io.NopCloser(bytes.NewBuffer(requestBodyBytes))

			body, _ := rest.Bind[any](c)

			request.Body = io.NopCloser(bytes.NewBuffer(requestBodyBytes))

			ctx, span := otel.Tracer("").Start(
				otelCtx,
				route,
				trace.WithSpanKind(trace.SpanKindServer),
				trace.WithAttributes(
					attribute.String("method", request.Method),
					attribute.String("path", path),
					attribute.String("route", route),
				),
			)

			span.AddEvent("request", trace.WithAttributes(
				attribute.String("path", path),
				attribute.String("method", request.Method),
				attribute.String("time", requestTime.Format(time.RFC3339)),
			))

			c.SetRequest(c.Request().WithContext(ctx))

			defer span.End()

			params := map[string]string{}

			for i, param := range c.ParamNames() {
				params[param] = c.ParamValues()[i]
			}

			userId := c.Get("user_id")

			resCapturer := &responseBodyCapturer{
				ResponseWriter: c.Response().Writer,
				body:           &bytes.Buffer{},
			}

			c.Response().Writer = resCapturer

			err := next(c)

			if err != nil {
				logging.Error(
					ctx,
					err,
				).Msg("error processing http request")
			}

			var responseBody any

			json.Unmarshal(resCapturer.body.Bytes(), &responseBody)

			requestLog := map[string]any{
				"path":    path,
				"method":  request.Method,
				"time":    requestTime.Format(time.RFC3339),
				"params":  params,
				"body":    body,
				"query":   c.QueryParams(),
				"userId":  userId,
				"headers": c.Request().Header,
			}

			responseLog := map[string]any{
				"code":        c.Response().Status,
				"processTime": time.Since(requestTime).Milliseconds(),
				"body":        responseBody,
				"headers":     c.Response().Header(),
			}

			span.AddEvent("response", trace.WithAttributes(
				attribute.Int("code", c.Response().Status),
				attribute.Int64("processTime", time.Since(requestTime).Milliseconds()),
			))

			if err != nil {
				responseLog["error"] = err.Error()
			}

			logging.Info(ctx).
				Interface("request", requestLog).
				Interface("response", responseLog).
				Msg("http.request")

			return err
		}
	}
}
