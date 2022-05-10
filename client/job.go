package client

import (
	"encoding/json"
	"errors"
	"sort"

	models "code.cloudfoundry.org/cli/plugin/models"
	"github.com/ess/hype"
	"github.com/starkandwayne/ocf-scheduler-cf-plugin/core"
	scheduler "github.com/starkandwayne/scheduler-for-ocf/core"
)

func ListJobs(driver *core.Driver, space models.SpaceFields) ([]*scheduler.Job, error) {
	params := hype.Params{}
	params.Set("space_guid", space.Guid)

	response := driver.Get("jobs", params)
	if !response.Okay() {
		return nil, response.Error()
	}

	data := struct {
		Resources []*scheduler.Job `json:"resources"`
	}{}

	err := json.Unmarshal(response.Data(), &data)
	if err != nil {
		return nil, err
	}

	return data.Resources, nil
}

type byExecutionStart []*scheduler.Execution

func (s byExecutionStart) Len() int {
	return len(s)
}

func (s byExecutionStart) Swap(i int, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byExecutionStart) Less(i int, j int) bool {
	return s[i].ExecutionStartTime.Before(s[j].ExecutionStartTime)
}

func ListJobExecutions(driver *core.Driver, job *scheduler.Job) ([]*scheduler.Execution, error) {
	response := driver.Get("jobs/"+job.GUID+"/history", nil)
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

	// Apparently, we only care about *scheduled* executions, not exeuctions
	// from ad-hock job runs.
	scheduled := make([]*scheduler.Execution, 0)
	for _, execution := range executions {
		if !execution.ScheduledTime.IsZero() {
			scheduled = append(scheduled, execution)
		}
	}

	sort.Sort(byExecutionStart(scheduled))

	return scheduled, nil
}

func ListJobSchedules(driver *core.Driver, job *scheduler.Job) ([]*scheduler.Schedule, error) {
	response := driver.Get("jobs/"+job.GUID+"/schedules", nil)
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

func JobNamed(driver *core.Driver, space models.SpaceFields, name string) (*scheduler.Job, error) {
	all, err := ListJobs(driver, space)
	if err != nil {
		return nil, err
	}

	for _, candidate := range all {
		if candidate.Name == name {
			return candidate, nil
		}
	}

	return nil, errors.New("no matching job found")
}

func CreateJob(driver *core.Driver, appGUID, name, command string) (*scheduler.Job, error) {
	params := hype.Params{}
	params.Set("app_guid", appGUID)

	input := &scheduler.Job{
		Name:    name,
		Command: command,
	}

	data, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	response := driver.Post("jobs", params, data)
	if !response.Okay() {
		return nil, response.Error()
	}

	output := &scheduler.Job{}
	err = json.Unmarshal(response.Data(), output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func DeleteJob(driver *core.Driver, job *scheduler.Job) error {
	response := driver.Delete("jobs/"+job.GUID, nil)
	if !response.Okay() {
		return response.Error()
	}

	return nil
}

func DeleteJobSchedule(driver *core.Driver, job *scheduler.Job, scheduleGUID string) error {
	response := driver.Delete(
		"jobs/"+job.GUID+"/schedules/"+scheduleGUID,
		nil,
	)

	if !response.Okay() {
		return response.Error()
	}

	return nil
}

func ExecuteJob(driver *core.Driver, job *scheduler.Job) (*scheduler.Execution, error) {
	response := driver.Post("jobs/"+job.GUID+"/execute", nil, nil)
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

func ScheduleJob(driver *core.Driver, job *scheduler.Job, expression string) (*scheduler.Schedule, error) {
	schedule := &scheduler.Schedule{
		Enabled:        true,
		Expression:     expression,
		ExpressionType: "cron_expression",
	}

	input, err := json.Marshal(schedule)
	if err != nil {
		return nil, err
	}

	response := driver.Post("jobs/"+job.GUID+"/schedules", nil, input)
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
