package main

import (
	"flag"
	"io"
	"testing"
)

func TestMainInvalidVersionType(t *testing.T) {
	flags := flag.NewFlagSet("bump_version", flag.ContinueOnError)
	flags.SetOutput(io.Discard)

	code := _main(flags, []string{"ptach", "server/serve.go"})
	if code != 2 {
		t.Fatalf("_main returned %d, want 2", code)
	}
}
