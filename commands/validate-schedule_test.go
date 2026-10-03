package commands

import (
	"strings"
	"testing"
	"time"
)

func TestParseValidateFrom(t *testing.T) {
	berlin, _ := time.LoadLocation("Europe/Berlin")
	for _, tt := range []struct {
		args []string
		want time.Time
	}{
		{[]string{"validate-schedule", "x", "--from", "2030-01-01"}, time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)},
		{[]string{"validate-schedule", "x", "--from", "2030-01-01T09:30"}, time.Date(2030, 1, 1, 9, 30, 0, 0, time.UTC)},
		{[]string{"validate-schedule", "x", "--from", "2030-01-01", "-t", "Europe/Berlin"}, time.Date(2030, 1, 1, 0, 0, 0, 0, berlin)},
		{[]string{"validate-schedule", "x", "--from", "2030-01-01T09:00:00-08:00", "-t", "Europe/Berlin"}, time.Date(2030, 1, 1, 17, 0, 0, 0, time.UTC)},
	} {
		opts, _, err := parseValidateFlags(tt.args)
		if err != nil || opts.from == nil || !opts.from.Equal(tt.want) {
			t.Errorf("%q: got %v %v, want %v", tt.args, opts.from, err, tt.want)
		}
	}
	if opts, _, err := parseValidateFlags([]string{"validate-schedule", "x"}); err != nil || opts.from != nil {
		t.Errorf("no --from: got %v %v", opts.from, err)
	}
}

func TestParseValidateFlags(t *testing.T) {
	opts, rest, err := parseValidateFlags([]string{"validate-schedule", "backup", "-t", "Europe/Berlin", "--prev", "2"})
	if err != nil || opts.next != nil || opts.prev != 2 || opts.timezone != "Europe/Berlin" || len(rest) != 1 || rest[0] != "backup" {
		t.Errorf("got %+v %q %v", opts, rest, err)
	}
	if opts, _, err := parseValidateFlags([]string{"validate-schedule", "x", "--next", "0"}); err != nil || opts.next == nil || *opts.next != 0 {
		t.Errorf("--next 0: got %+v %v", opts, err)
	}
	for _, tt := range []struct {
		args []string
		want string
	}{
		{[]string{"validate-schedule", "x", "--next", "-5"}, "--next must be 0 or more"},
		{[]string{"validate-schedule", "x", "--next", "-1"}, "--next must be 0 or more"},
		{[]string{"validate-schedule", "x", "-t", "Europe/ Berlin"}, "no spaces allowed in the timezone"},
		{[]string{"validate-schedule", "x", "--from", "next tuesday"}, "--from must be"},
		{[]string{"validate-schedule", "x", "--from", "2030-01-01", "-t", "Mars/Olympus"}, "unknown time zone"},
	} {
		if _, _, err := parseValidateFlags(tt.args); err == nil || !strings.Contains(err.Error(), tt.want) {
			t.Errorf("%q: got %v, want %q", tt.args, err, tt.want)
		}
	}
}
