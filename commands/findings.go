package commands

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	scheduler "github.com/cloudfoundry-community/ocf-scheduler/core"

	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/client"
)

// renderFindings prints errors, with a caret under the offending value when
// the scheduler located it.
func renderFindings(w io.Writer, expression string, findings []scheduler.Finding) {
	for _, f := range findings {
		fmt.Fprintln(w, f.Message)
		if f.Offset != nil && *f.Offset <= len(expression) {
			fmt.Fprintln(w, "  "+expression)
			fmt.Fprintln(w, "  "+strings.Repeat(" ", *f.Offset)+strings.Repeat("^", max(1, len(f.Value))))
		}
	}
}

// printRejection shows why the scheduler refused a schedule.
func printRejection(err error, expression string) {
	var rejected *client.RejectedError
	if errors.As(err, &rejected) && len(rejected.Findings.Errors) > 0 {
		renderFindings(os.Stdout, expression, rejected.Findings.Errors)
		return
	}
	fmt.Println(err.Error())
}
