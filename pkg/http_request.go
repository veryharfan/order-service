package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/app/handler/response"
	"order-service/pkg/ctxutil"

	"github.com/mitchellh/mapstructure"
)

type AuthInternalHeader string

const (
	AuthInternalHeaderKey AuthInternalHeader = "X-Internal-Auth"
)

func AddRequestHeader(ctx context.Context, internalAuthHeader string, httpRequest *http.Request) {
	httpRequest.Header.Add("Content-Type", "application/json")
	httpRequest.Header.Add("Accept", "application/json")

	if reqID := ctx.Value(ctxutil.RequestIDKey); reqID != nil {
		httpRequest.Header.Add("X-Request-ID", reqID.(string))
	}

	httpRequest.Header.Add(string(AuthInternalHeaderKey), internalAuthHeader)
}

func DecodeResponseBody[T any](resp *http.Response, v T) error {
	var respBody response.Response[T]
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return fmt.Errorf("failed to decode response body: %w", err)
	}

	if !respBody.Success {
		return fmt.Errorf("response error: %s", respBody.Error)
	}

	err := mapstructure.Decode(respBody.Data, &v)
	if err != nil {
		return fmt.Errorf("failed to map response data: %w", err)
	}

	return nil
}
