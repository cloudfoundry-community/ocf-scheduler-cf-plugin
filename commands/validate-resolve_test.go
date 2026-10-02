package commands

import (
	"errors"
	"strings"
	"testing"

	scheduler "github.com/cloudfoundry-community/ocf-scheduler/core"
)

var (
	jobs = []*scheduler.Job{
		{GUID: "j1", Name: "backup", AppGUID: "a1"},
		{GUID: "j2", Name: "nightly report", AppGUID: "a1"},
		{GUID: "j3", Name: "dup", AppGUID: "a1"},
		{GUID: "j4", Name: "dup", AppGUID: "a2"},
		{GUID: "j5", Name: "both", AppGUID: "a1"},
	}
	calls = []*scheduler.Call{
		{GUID: "c1", Name: "ping", AppGUID: "a1"},
		{GUID: "c2", Name: "both", AppGUID: "a1"},
	}
)

func TestResolveArgs(t *testing.T) {
	tests := []struct {
		args       []string
		ref        string // "type/guid" or "" for a bare expression
		expression string
		err        string
	}{
		{[]string{"backup"}, "job/j1", "", ""},
		{[]string{"backup", "H 2 * * *"}, "job/j1", "H 2 * * *", ""},
		{[]string{"j1", "0 2 * * *"}, "job/j1", "0 2 * * *", ""},
		{[]string{"nightly report"}, "job/j2", "", ""},
		{[]string{"ping"}, "call/c1", "", ""},
		{[]string{"0 2 * * *"}, "", "0 2 * * *", ""},
		{[]string{"@daily"}, "", "@daily", ""},
		{[]string{"job", "both", "0 2 * * *"}, "job/j5", "0 2 * * *", ""},
		{[]string{"call", "both"}, "call/c2", "", ""},
		{[]string{"job"}, "", "job", ""}, // a lone word is an expression unless something is named "job"
		{[]string{"both"}, "", "", "matches 2 jobs or calls"},
		{[]string{"dup"}, "", "", "matches 2 jobs or calls"},
		{[]string{"call", "backup"}, "", "", `no call named "backup"`},
		{[]string{"0", "2", "*", "*", "*"}, "", "", "expected [job|call]"},
		{[]string{"0", "2 * * *"}, "", "", "quote a cron expression"},
		{nil, "", "", "expected [job|call]"},
	}
	for _, tt := range tests {
		got, err := resolveArgs(tt.args, jobs, calls)
		if tt.err != "" {
			if err == nil || !strings.Contains(err.Error(), tt.err) {
				t.Errorf("%q: error %v, want %q", tt.args, err, tt.err)
			}
			continue
		}
		ref := ""
		if got.Target != nil {
			ref = got.Target.Type + "/" + got.Target.GUID
		}
		if err != nil || ref != tt.ref || got.Expression != tt.expression {
			t.Errorf("%q: got %s %q, %v", tt.args, ref, got.Expression, err)
		}
	}
	var amb *AmbiguousError
	if _, err := resolveArgs([]string{"dup"}, jobs, calls); !errors.As(err, &amb) || len(amb.Matches) != 2 || amb.Matches[1].AppGUID != "a2" {
		t.Errorf("ambiguous: %v", err)
	}
}
