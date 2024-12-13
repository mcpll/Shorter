package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGenerateShortUrl(t *testing.T) {
	tests := []struct {
		in   string
		out  string
		port string
	}{
		{in: "http://www.google.it", out: "http://localhost:8080/short/yQ9aVE", port: "8080"},
		{in: "http://www.google.it", out: "http://localhost:2020/short/yQ9aVE", port: "2020"},
	}

	var db ShortUrls
	db.urls = make(map[string]string)

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			json := fmt.Sprintf(`{"url":"%s"}`, test.in)
			address := fmt.Sprintf("http://localhost:%s/create", test.port)
			req := httptest.NewRequest("POST", address, strings.NewReader(json))
			w := httptest.NewRecorder()
			handler := handleGenerateShortUrl(&db, test.port)

			handler(w, req)
			resp := w.Result()

			body, _ := io.ReadAll(resp.Body)

			want := test.out
			got := string(body)

			if got != want {
				t.Fatalf("got: %s - want: %s", got, want)
				return
			}
		})
	}
}

func TestRedirect(t *testing.T) {
	tests := []struct {
		shortCode      string
		url            string
		port           string
		redirectStatus int
	}{
		{shortCode: "yQ9aVE", url: "http://www.google.it", port: "8080", redirectStatus: http.StatusMovedPermanently},
	}

	for _, test := range tests {
		t.Run(test.shortCode, func(t *testing.T) {
			db := ShortUrls{urls: map[string]string{test.shortCode: test.url}}
			address := fmt.Sprintf("http://localhost:%s/short/%s", test.port, test.shortCode)

			req := httptest.NewRequest("GET", address, nil)
			w := httptest.NewRecorder()

			handler := handleRedirect(&db)
			handler(w, req)

			resp := w.Result()

			want := test.redirectStatus
			got := resp.StatusCode

			if want != got {
				t.Fatalf("got: %v - want: %v", got, want)
				return
			}
		})
	}
}
