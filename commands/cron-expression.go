package commands

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/pflag"
	"golang.org/x/term"
)

// cf cron-expression
func CronExpression(content string, args []string) {
	flags := pflag.NewFlagSet("cron-expression", pflag.ExitOnError)
	noPager := flags.BoolP("no-pager", "n", false, "Do not pipe output through a pager")
	flags.Parse(args)

	if *noPager || !term.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Print(content)
		return
	}

	if err := runPager(content); err != nil {
		fmt.Print(content)
	}
}

func runPager(content string) error {
	pager := os.Getenv("PAGER")

	var cmd *exec.Cmd
	if pager != "" {
		cmd = exec.Command(pager)
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
