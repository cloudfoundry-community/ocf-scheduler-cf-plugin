package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"

	"github.com/spf13/pflag"
	"golang.org/x/term"
)

// cf cron-expression
func CronExpression(styles map[string]string, args []string) {
	flags := pflag.NewFlagSet("cron-expression", pflag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	noPager := flags.BoolP("no-pager", "n", false, "Do not pipe output through a pager")
	style := flags.StringP("style", "s", "", "Glamour style for rendering")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, pflag.ErrHelp) {
			fmt.Fprintln(os.Stderr, "\nFor full usage details run: cf help cron-expression")
		}
		os.Exit(1)
	}

	validStyles := AvailableStyles(styles)
	if len(validStyles) == 0 {
		fmt.Fprintln(os.Stderr, "Error: no glamour styles available")
		os.Exit(1)
	}

	if !flags.Changed("style") {
		if term.IsTerminal(int(os.Stdout.Fd())) {
			*style = "dark"
		} else {
			*style = "notty"
		}
		// Fall back to first available style if the default is missing.
		if _, ok := styles[*style]; !ok {
			fmt.Fprintf(os.Stderr, "Warning: default style %q unavailable, using %q\n", *style, validStyles[0])
			*style = validStyles[0]
		}
	}

	content, ok := styles[*style]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown style %q. Valid styles: %s\n", *style, strings.Join(validStyles, ", "))
		os.Exit(1)
	}

	if *noPager || !term.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Print(content)
		return
	}

	if err := runPager(content); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: pager failed (%v), printing directly\n", err)
		if fallback, ok := styles["notty"]; ok {
			fmt.Print(fallback)
		} else {
			// Use first available style as a last resort.
			fmt.Print(styles[validStyles[0]])
		}
	}
}

func AvailableStyles(styles map[string]string) []string {
	names := make([]string, 0, len(styles))
	for name := range styles {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func runPager(content string) error {
	pager := os.Getenv("PAGER")

	var cmd *exec.Cmd
	if pager != "" {
		parts := strings.Fields(pager)
		cmd = exec.Command(parts[0], parts[1:]...)
	} else if runtime.GOOS == "windows" {
		cmd = exec.Command("more")
	} else {
		cmd = exec.Command("less", "-FRX")
	}

	cmd.Stdin = strings.NewReader(content)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
