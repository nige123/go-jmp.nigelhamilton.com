package main

import (
    "os/exec"
    "strings"
    "testing"

    "github.com/nige123/go-jmp.nigelhamilton.com/pkg/version"
)

func TestVersionCommand(t *testing.T) {
    cmd := exec.Command("go", "run", ".", "version")
    cmd.Dir = "."

    output, err := cmd.CombinedOutput()
    if err != nil {
        t.Fatalf("expected version command to succeed, got error: %v, output: %s", err, string(output))
    }

    expected := "jmp - version " + version.VERSION
    got := strings.TrimSpace(string(output))
    if got != expected {
        t.Fatalf("expected %q, got %q", expected, got)
    }
}
