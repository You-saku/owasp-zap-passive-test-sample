package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuth(t *testing.T) {
	srv := httptest.NewServer(newMux(nil, "secret"))
	defer srv.Close()

	for _, tc := range []struct {
		header string
		want   int
	}{
		{"", http.StatusUnauthorized},
		{"Bearer wrong", http.StatusUnauthorized},
		{"Bearer secret", http.StatusOK},
	} {
		req, _ := http.NewRequest("GET", srv.URL+"/me", nil)
		req.Header.Set("Authorization", tc.header)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != tc.want {
			t.Errorf("header %q: got %d, want %d", tc.header, res.StatusCode, tc.want)
		}
	}
}
