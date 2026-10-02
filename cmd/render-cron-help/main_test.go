package main

import (
	"strings"
	"testing"
)

func TestRenderWrapsHyphenatedParagraph(t *testing.T) {
	src := "Expressions are parsed by the go-cron library, which supports standard cron syntax, optional seconds and year fields, extended day-of-month/day-of-week syntax, descriptor shortcuts, and timezone prefixes.\n"
	want := []string{
		"",
		"  Expressions are parsed by the go-cron library, which supports standard",
		"  cron syntax, optional seconds and year fields, extended",
		"  day-of-month/day-of-week syntax, descriptor shortcuts, and timezone",
		"  prefixes.",
		"",
		"",
	}

	out, err := render(src, "notty")
	if err != nil {
		t.Fatal(err)
	}
	got := strings.Split(out, "\n")
	for i := range got {
		got[i] = strings.TrimRight(got[i], " ")
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("got:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

func TestRenderKeepsTables(t *testing.T) {
	src := "| Field | Range |\n|---|---|\n| Minute | `0-59` |\n"

	out, err := render(src, "notty")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "|---|") || !strings.Contains(out, "0-59") {
		t.Errorf("table not rendered:\n%s", out)
	}
}
