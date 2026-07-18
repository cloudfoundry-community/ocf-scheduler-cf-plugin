package core

import (
	"fmt"
	"strings"
	"testing"
)

type fakeCurler struct {
	responses map[string][]string
	requests  []string
	err       error
}

func (f *fakeCurler) CliCommandWithoutTerminalOutput(args ...string) ([]string, error) {
	f.requests = append(f.requests, strings.Join(args, " "))
	if f.err != nil {
		return nil, f.err
	}
	if len(args) != 2 || args[0] != "curl" {
		return nil, fmt.Errorf("unexpected command: %v", args)
	}
	resp, ok := f.responses[args[1]]
	if !ok {
		return nil, fmt.Errorf("no stubbed response for %s", args[1])
	}
	return resp, nil
}

func TestAppsV3SinglePage(t *testing.T) {
	curler := &fakeCurler{responses: map[string][]string{
		"/v3/apps?space_guids=space-1&per_page=5000": {
			`{"pagination":{"next":null},"resources":[`,
			`{"guid":"app-guid-1","name":"app-one"},`,
			`{"guid":"app-guid-2","name":"app-two"}]}`,
		},
	}}

	apps, err := AppsV3(curler, "space-1")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(apps) != 2 {
		t.Fatalf("expected 2 apps, got %d", len(apps))
	}
	if apps[0].Guid != "app-guid-1" || apps[0].Name != "app-one" {
		t.Errorf("unexpected first app: %+v", apps[0])
	}
	if apps[1].Guid != "app-guid-2" || apps[1].Name != "app-two" {
		t.Errorf("unexpected second app: %+v", apps[1])
	}
}

func TestAppsV3FollowsPagination(t *testing.T) {
	curler := &fakeCurler{responses: map[string][]string{
		"/v3/apps?space_guids=space-1&per_page=5000": {
			`{"pagination":{"next":{"href":"https://api.example.com/v3/apps?page=2&per_page=5000&space_guids=space-1"}},`,
			`"resources":[{"guid":"g1","name":"a1"}]}`,
		},
		"/v3/apps?page=2&per_page=5000&space_guids=space-1": {
			`{"pagination":{"next":null},"resources":[{"guid":"g2","name":"a2"}]}`,
		},
	}}

	apps, err := AppsV3(curler, "space-1")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(apps) != 2 {
		t.Fatalf("expected 2 apps across pages, got %d", len(apps))
	}
	if apps[1].Guid != "g2" {
		t.Errorf("expected paginated app g2, got %+v", apps[1])
	}
}

func TestAppV3ByNameFound(t *testing.T) {
	curler := &fakeCurler{responses: map[string][]string{
		"/v3/apps?names=app-one&space_guids=space-1&per_page=5000": {
			`{"pagination":{"next":null},"resources":[{"guid":"app-guid-1","name":"app-one"}]}`,
		},
	}}

	app, err := AppV3ByName(curler, "space-1", "app-one")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if app.Guid != "app-guid-1" {
		t.Errorf("unexpected app: %+v", app)
	}
}

func TestAppV3ByNameEscapesQuery(t *testing.T) {
	curler := &fakeCurler{responses: map[string][]string{
		"/v3/apps?names=app+one%26x&space_guids=space-1&per_page=5000": {
			`{"pagination":{"next":null},"resources":[{"guid":"g","name":"app one&x"}]}`,
		},
	}}

	app, err := AppV3ByName(curler, "space-1", "app one&x")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if app.Guid != "g" {
		t.Errorf("unexpected app: %+v", app)
	}
}

func TestAppV3ByNameNotFound(t *testing.T) {
	curler := &fakeCurler{responses: map[string][]string{
		"/v3/apps?names=ghost&space_guids=space-1&per_page=5000": {
			`{"pagination":{"next":null},"resources":[]}`,
		},
	}}

	_, err := AppV3ByName(curler, "space-1", "ghost")
	if err == nil {
		t.Fatal("expected an error for a missing app")
	}
}

func TestAppsV3CCErrorSurfaces(t *testing.T) {
	curler := &fakeCurler{responses: map[string][]string{
		"/v3/apps?space_guids=space-1&per_page=5000": {
			`{"errors":[{"code":10008,"title":"CF-UnprocessableEntity","detail":"something broke"}]}`,
		},
	}}

	_, err := AppsV3(curler, "space-1")
	if err == nil {
		t.Fatal("expected an error when CC returns an errors payload")
	}
	if !strings.Contains(err.Error(), "something broke") {
		t.Errorf("expected CC error detail in message, got: %s", err)
	}
}
