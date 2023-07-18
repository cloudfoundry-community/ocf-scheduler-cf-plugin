package client

import (
	"encoding/json"
	"errors"
	"sort"

	models "code.cloudfoundry.org/cli/plugin/models"
	"github.com/ess/hype"
	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/core"
	scheduler "github.com/cloudfoundry-community/ocf-scheduler/core"
)

func ListCalls(driver *core.Driver, space models.SpaceFields) ([]*scheduler.Call, error) {
	params := hype.Params{}
	params.Set("space_guid", space.Guid)

	response := driver.Get("calls", params)
	if !response.Okay() {
		return nil, response.Error()
	}

	data := struct {
		Resources []*scheduler.Call `json:"resources"`
	}{}

	err := json.Unmarshal(response.Data(), &data)
	if err != nil {
		return nil, err
	}

	return data.Resources, nil
}

func ListCallExecutions(driver *core.Driver, call *scheduler.Call) ([]*scheduler.Execution, error) {
	response := driver.Get("calls/"+call.GUID+"/history", nil)
	if !response.Okay() {
		return nil, response.Error()
	}

	data := struct {
		Resources []*scheduler.Execution `json:"resources"`
	}{}

	err := json.Unmarshal(response.Data(), &data)
	if err != nil {
		return nil, err
	}

	executions := data.Resources

	// Apparently, we only care about *scheduled* executions, not executions
	// from ad-hoc call runs.
	scheduled := make([]*scheduler.Execution, 0)
	for _, execution := range executions {
		if !execution.ScheduledTime.IsZero() {
			scheduled = append(scheduled, execution)
		}
	}

	sort.Sort(byExecutionStart(scheduled))

	return scheduled, nil
}

func ListCallSchedules(driver *core.Driver, call *scheduler.Call) ([]*scheduler.Schedule, error) {
	response := driver.Get("calls/"+call.GUID+"/schedules", nil)

	if !response.Okay() {
		return nil, response.Error()
	}

	data := struct {
		Resources []*scheduler.Schedule `json:"resources"`
	}{}

	err := json.Unmarshal(response.Data(), &data)
	if err != nil {
		return nil, err
	}

	return data.Resources, nil
}

func GetCall(driver *core.Driver, callGUID string) (*scheduler.Call, error) {
	response := driver.Get("calls/"+callGUID, nil)
	if !response.Okay() {
		return nil, response.Error()
	}

	call := &scheduler.Call{}
	err := json.Unmarshal(response.Data(), call)
	if err != nil {
		return nil, err
	}

	return call, nil
}

func CallNamed(driver *core.Driver, space models.SpaceFields, name string) (*scheduler.Call, error) {
	all, err := ListCalls(driver, space)
	if err != nil {
		return nil, err
	}

	for _, candidate := range all {
		if candidate.Name == name {
			return candidate, nil
		}
	}

	return nil, errors.New("no matching call found")
}

func CreateCall(driver *core.Driver, appGUID, name, url string) (*scheduler.Call, error) {
	params := hype.Params{}
	params.Set("app_guid", appGUID)

	input := &scheduler.Call{
		Name: name,
		URL:  url,
		// TODO: Figure out what this should actually be.
		AuthHeader: "default",
	}

	data, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	response := driver.Post("calls", params, data)
	if !response.Okay() {
		return nil, response.Error()
	}

	output := &scheduler.Call{}
	err = json.Unmarshal(response.Data(), output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func DeleteCall(driver *core.Driver, call *scheduler.Call) error {
	response := driver.Delete("calls/"+call.GUID, nil)
	if !response.Okay() {
		return response.Error()
	}

	return nil
}

func DeleteCallSchedule(driver *core.Driver, call *scheduler.Call, scheduleGUID string) error {
	response := driver.Delete(
		"calls/"+call.GUID+"/schedules/"+scheduleGUID,
		nil,
	)

	if !response.Okay() {
		return response.Error()
	}

	return nil
}

func ExecuteCall(driver *core.Driver, call *scheduler.Call) (*scheduler.Execution, error) {
	response := driver.Post("calls/"+call.GUID+"/execute", nil, nil)
	if !response.Okay() {
		return nil, response.Error()
	}

	data := &scheduler.Execution{}
	err := json.Unmarshal(response.Data(), data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func ScheduleCall(driver *core.Driver, call *scheduler.Call, expression string) (*scheduler.Schedule, error) {
	schedule := &scheduler.Schedule{
		Enabled:        true,
		Expression:     expression,
		ExpressionType: "cron_expression",
	}

	input, err := json.Marshal(schedule)
	if err != nil {
		return nil, err
	}

	response := driver.Post("calls/"+call.GUID+"/schedules", nil, input)
	if !response.Okay() {
		return nil, response.Error()
	}

	output := &scheduler.Schedule{}
	err = json.Unmarshal(response.Data(), output)
	if err != nil {
		return nil, err
	}

	return output, nil
}
