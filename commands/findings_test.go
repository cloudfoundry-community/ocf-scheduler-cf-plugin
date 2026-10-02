package commands

import (
	"bytes"
	"testing"

	scheduler "github.com/cloudfoundry-community/ocf-scheduler/core"
)

func TestRenderFindings(t *testing.T) {
	off := 12
	var b bytes.Buffer
	renderFindings(&b, "0 0 * * MON-FUN", []scheduler.Finding{
		{Code: "parse_error", Message: "day-of-week field value FUN is not a number or name (valid: 0-7)", Value: "FUN", Offset: &off},
		{Code: "never_fires", Message: "this schedule never runs"},
	})
	want := "day-of-week field value FUN is not a number or name (valid: 0-7)\n  0 0 * * MON-FUN\n              ^^^\nthis schedule never runs\n"
	if b.String() != want {
		t.Errorf("got\n%q\nwant\n%q", b.String(), want)
	}
}
