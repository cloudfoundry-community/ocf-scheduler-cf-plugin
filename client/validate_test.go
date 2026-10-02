package client

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/core"
	scheduler "github.com/cloudfoundry-community/ocf-scheduler/core"
)

func serve(t *testing.T, status int, body string) *core.Driver {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)
	driver, err := core.NewDriver(server.URL, "Bearer token")
	if err != nil {
		t.Fatal(err)
	}
	return driver
}

func TestValidateExpression(t *testing.T) {
	driver := serve(t, http.StatusOK, `{"expression":"*/15 * * * *","valid":true,"description":"every 15 minutes"}`)
	got, err := ValidateExpression(driver, scheduler.ValidateRequest{Expression: "*/15 * * * *"})
	if err != nil || !got.Valid || got.Description != "every 15 minutes" {
		t.Errorf("got %+v, %v", got, err)
	}
}

func TestValidateStored(t *testing.T) {
	driver := serve(t, http.StatusOK, `{"ref":{"type":"job","guid":"g","name":"backup"},"resources":[{"schedule_guid":"s1","valid":true}]}`)
	got, err := ValidateStored(driver, scheduler.ValidateRequest{RefType: "job", RefGUID: "g"})
	if err != nil || got.Ref.Name != "backup" || len(got.Resources) != 1 || got.Resources[0].ScheduleGUID != "s1" {
		t.Errorf("got %+v, %v", got, err)
	}
}

func TestValidateOldScheduler(t *testing.T) {
	driver := serve(t, http.StatusNotFound, `{"message":"Not Found"}`)
	if _, err := ValidateExpression(driver, scheduler.ValidateRequest{Expression: "x"}); !errors.Is(err, ErrValidationUnsupported) {
		t.Errorf("got %v", err)
	}
}

func TestValidateRejected(t *testing.T) {
	for _, tt := range []struct {
		status int
		body   string
		msg    string
	}{
		{http.StatusNotFound, `{"errors":[{"code":"ref_not_found","message":"job g not found"}],"warnings":[]}`, "job g not found"},
		{http.StatusUnprocessableEntity, `{"errors":[{"code":"bad_request","message":"next and prev must be 0 to 100"}],"warnings":[]}`, "next and prev must be 0 to 100"},
		{http.StatusUnauthorized, `""`, "response status: 401"},
		{http.StatusUnprocessableEntity, `"expected exactly 5 fields, found 6: [0 0 0 * * *]"`, "expected exactly 5 fields, found 6: [0 0 0 * * *]"},
	} {
		_, err := ValidateExpression(serve(t, tt.status, tt.body), scheduler.ValidateRequest{Expression: "x"})
		var rejected *RejectedError
		if !errors.As(err, &rejected) || err.Error() != tt.msg {
			t.Errorf("%d: got %v", tt.status, err)
		}
	}
}
