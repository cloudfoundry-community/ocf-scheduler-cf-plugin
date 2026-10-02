package client

import (
	"errors"
	"net/http"
	"testing"

	scheduler "github.com/cloudfoundry-community/ocf-scheduler/core"
)

func TestScheduleJobRejected(t *testing.T) {
	driver := serve(t, http.StatusUnprocessableEntity, `{"errors":[{"code":"parse_error","message":"minute field value 60 is out of range (valid: 0-59)","field":"minute","value":"60","offset":0}],"warnings":[]}`)
	_, err := ScheduleJob(driver, &scheduler.Job{GUID: "g"}, "60 * * * *")
	var rejected *RejectedError
	if !errors.As(err, &rejected) || rejected.Findings.Errors[0].Field != "minute" || *rejected.Findings.Errors[0].Offset != 0 {
		t.Errorf("got %v", err)
	}
	driver = serve(t, http.StatusCreated, `{"guid":"s1","expression":"0 2 * * *"}`)
	if s, err := ScheduleCall(driver, &scheduler.Call{GUID: "g"}, "0 2 * * *"); err != nil || s.GUID != "s1" {
		t.Errorf("call: got %+v, %v", s, err)
	}
}
