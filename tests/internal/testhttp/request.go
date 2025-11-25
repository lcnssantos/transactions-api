package testhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type Request[T any] struct {
	Endpoint string
	Method   string
	Body     *T
}

type Response[T any] struct {
	Code int
	Body *T
}

func DoRequest[Req any, Res any](ctx context.Context, request Request[Req]) (Response[Res], error) {
	var body io.Reader

	if request.Body != nil {
		payload, err := json.Marshal(request.Body)

		if err != nil {
			return Response[Res]{}, err
		}

		body = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(ctx, request.Method, request.Endpoint, body)

	if err != nil {
		return Response[Res]{}, err
	}

	if request.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return Response[Res]{}, err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return Response[Res]{}, err
	}

	var data *Res

	if len(respBody) > 0 {
		var parsed Res

		if err := json.Unmarshal(respBody, &parsed); err != nil {
			return Response[Res]{}, err
		}

		data = &parsed
	}

	return Response[Res]{
		Code: resp.StatusCode,
		Body: data,
	}, nil
}
