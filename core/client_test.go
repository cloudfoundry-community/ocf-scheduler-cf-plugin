package core

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostJSONKeepsErrorBodies(t *testing.T) {
	var got *http.Request
	var sent string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r
		b, _ := io.ReadAll(r.Body)
		sent = string(b)
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{"errors":[{"code":"parse_error"}]}`))
	}))
	defer server.Close()

	driver, err := NewDriver(server.URL+"/", "Bearer token")
	if err != nil {
		t.Fatal(err)
	}
	status, body, err := driver.PostJSON("schedules/validate", []byte(`{"expression":"x"}`))
	if err != nil || status != http.StatusUnprocessableEntity || string(body) != `{"errors":[{"code":"parse_error"}]}` {
		t.Errorf("got %d %q %v", status, body, err)
	}
	if got.URL.Path != "/schedules/validate" || got.Header.Get("Authorization") != "Bearer token" ||
		got.Header.Get("Content-Type") != "application/json" || sent != `{"expression":"x"}` {
		t.Errorf("request %s %v body %q", got.URL.Path, got.Header, sent)
	}
}
