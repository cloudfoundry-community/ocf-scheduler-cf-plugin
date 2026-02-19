// render-cron-help reads cron-expression-reference.md and renders it
// with glamour's notty style to produce clean, formatted plain text
// suitable for terminal display without ANSI escape codes.
//
// Usage: go run ./cmd/render-cron-help
package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/glamour"
)

func main() {
	src, err := os.ReadFile("cron-expression-reference.md")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading markdown: %v\n", err)
		os.Exit(1)
	}

	tableStyle := []byte(`{"table":{"center_separator":"┼","column_separator":"│","row_separator":"─"}}`)

	renderer, err := glamour.NewTermRenderer(
		glamour.WithStylePath("notty"),
		glamour.WithStylesFromJSONBytes(tableStyle),
		glamour.WithWordWrap(72),
		glamour.WithTableWrap(true),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating renderer: %v\n", err)
		os.Exit(1)
	}

	rendered, err := renderer.Render(string(src))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering markdown: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile("cron-expression-reference.rendered", []byte(rendered), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing rendered file: %v\n", err)
		os.Exit(1)
	}
}
