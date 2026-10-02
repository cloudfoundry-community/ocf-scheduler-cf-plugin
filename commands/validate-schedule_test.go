package commands

import (
	"strings"
	"testing"
)

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
	} {
		if _, _, err := parseValidateFlags(tt.args); err == nil || !strings.Contains(err.Error(), tt.want) {
			t.Errorf("%q: got %v, want %q", tt.args, err, tt.want)
		}
	}
}
