package pp

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestPipeIsNotTerminal(t *testing.T) {
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		reader.Close()
		writer.Close()
	})

	if isTerminal(writer) {
		t.Fatal("pipe detected as a terminal")
	}
}

func TestPrintToDisablesColorForPipedOutput(t *testing.T) {
	var out bytes.Buffer
	printTo(&out, false, struct{ Subject string }{Subject: "example.com"})

	got := out.String()
	if strings.Contains(got, "\x1b[") {
		t.Fatalf("piped output contains ANSI escape sequences: %q", got)
	}
	if !strings.Contains(got, "Subject:") {
		t.Fatalf("piped output cannot be matched by field name: %q", got)
	}
}

func TestPrintToKeepsColorForTerminalOutput(t *testing.T) {
	var out bytes.Buffer
	printTo(&out, true, struct{ Subject string }{Subject: "example.com"})

	if !strings.Contains(out.String(), "\x1b[") {
		t.Fatalf("terminal output does not contain ANSI escape sequences: %q", out.String())
	}
}
