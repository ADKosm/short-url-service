package handlers

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"mainmod/core/gateways/redismock"
	"mainmod/core/handlers/generatormock"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_HandleAdd(t *testing.T) {
	redisClient := redismock.NewGateway()
	gen := generatormock.NewGenerator()

	genValue := "dfhsadqe"
	gen.GenShortURLFunc = func() string {
		return genValue
	}

	writeCnt := 0
	redisClient.WriteFunc = func(ctx context.Context, key string, value string) error {
		require.Equal(t, genValue, key)
		require.Equal(t, "http://yandex.ru", value)
		writeCnt++

		return nil
	}

	handler := NewHandler(redisClient, gen)

	req := httptest.NewRequest("POST", "/api/add", strings.NewReader(`
	{
		"url": "http://yandex.ru"
	}
    `))
	w := httptest.NewRecorder()
	handler.HandleAdd(w, req)

	result := w.Result()

	var resp AddResponse
	err := json.NewDecoder(result.Body).Decode(&resp)

	require.NoError(t, err)
	require.Equal(t, "dfhsadqe", resp.ShortUrl)
	require.Equal(t, 1, writeCnt)
}
