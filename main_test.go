package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		count int
		want  int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, len(cafeList["moscow"])},
	}

	for _, v := range requests {
		reqURL := fmt.Sprintf("/cafe?count=%d&city=moscow", v.count)

		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", reqURL, nil)

		handler.ServeHTTP(response, req)
		responseStr := strings.TrimSpace(response.Body.String())

		result := strings.Split(responseStr, ",")
		count := len(result)
		if responseStr == "" {
			count = 0
		}

		require.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, v.want, count)
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		search    string
		wantCount int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}

	for _, v := range requests {
		reqURL := fmt.Sprintf("/cafe?search=%s&city=moscow", v.search)

		response := httptest.NewRecorder()

		req := httptest.NewRequest("GET", reqURL, nil)

		handler.ServeHTTP(response, req)

		respStr := strings.TrimSpace(response.Body.String())
		result := strings.Split(respStr, ",")

		count := len(result)
		if respStr == "" {
			count = 0
		}

		assert.Equal(t, v.wantCount, count)

		if count > 0 {
			for _, cafeName := range result {
				cafeLower := strings.ToLower(cafeName)
				searchLower := strings.ToLower(v.search)

				check := strings.Contains(cafeLower, searchLower)

				assert.True(t, check)
			}
		}
	}
}
