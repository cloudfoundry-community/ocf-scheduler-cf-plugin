package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/core"
	scheduler "github.com/cloudfoundry-community/ocf-scheduler/core"
)

// MinValidateVersion is the first scheduler release with POST /schedules/validate.
const MinValidateVersion = "2.1.0"

// ErrValidationUnsupported means the scheduler predates POST /schedules/validate.
var ErrValidationUnsupported = errors.New("this scheduler does not support validate-schedule; it needs ocf-scheduler " + MinValidateVersion + " or later")

// RejectedError is a 4xx response that explains itself.
type RejectedError struct {
	Status   int
	Findings scheduler.Findings
}

func (e *RejectedError) Error() string {
	if len(e.Findings.Errors) == 0 {
		return fmt.Sprintf("response status: %d", e.Status)
	}
	return e.Findings.Errors[0].Message
}

// StoredValidation is the re-check of a job's or call's stored schedules.
type StoredValidation struct {
	Ref       *scheduler.Ref                `json:"ref"`
	Resources []*scheduler.ScheduleAnalysis `json:"resources"`
}

// ValidateExpression validates req.Expression, alone or for req's job or call.
func ValidateExpression(driver *core.Driver, req scheduler.ValidateRequest) (*scheduler.ScheduleAnalysis, error) {
	out := &scheduler.ScheduleAnalysis{}
	return out, validate(driver, req, out)
}

// ValidateStored re-checks the stored schedules of req's job or call.
func ValidateStored(driver *core.Driver, req scheduler.ValidateRequest) (*StoredValidation, error) {
	out := &StoredValidation{}
	return out, validate(driver, req, out)
}

func validate(driver *core.Driver, req scheduler.ValidateRequest, out any) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	status, data, err := driver.PostJSON("schedules/validate", body)
	if err != nil {
		return err
	}
	switch {
	case status == http.StatusOK:
		return json.Unmarshal(data, out)
	case status == http.StatusNotFound && !bytes.Contains(data, []byte(`"ref_not_found"`)):
		return ErrValidationUnsupported
	}
	return rejected(status, data)
}

// rejected decodes a 4xx body of {errors, warnings}; any other body is
// reported by status alone.
func rejected(status int, data []byte) error {
	e := &RejectedError{Status: status}
	if json.Unmarshal(data, &e.Findings) != nil {
		e.Findings = scheduler.Findings{}
	}
	return e
}
