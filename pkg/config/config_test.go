package config

import (
    "os"
    "path/filepath"
    "testing"
)

func TestDefaultFindTemplateUsesRipgrep(t *testing.T) {
    home := t.TempDir()
    t.Setenv("HOME", home)

    cfg, err := New()
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }

    template, err := cfg.Get("find.command.template")
    if err != nil {
        t.Fatalf("expected template value, got %v", err)
    }

    expected := "rg --line-number --with-filename --no-heading --color never '[-search-terms-]'"
    if template != expected {
        t.Fatalf("expected %q, got %q", expected, template)
    }

    rendered, err := cfg.GetTemplate("find.command.template", map[string]string{"search-terms": "needle"})
    if err != nil {
        t.Fatalf("expected rendered template, got %v", err)
    }

    checks := []string{"rg ", "--line-number", "--with-filename", "--no-heading", "--color never", "needle"}
    for _, check := range checks {
        if !contains(rendered, check) {
            t.Fatalf("expected rendered command to contain %q, got %q", check, rendered)
        }
    }

    if _, err := os.Stat(filepath.Join(home, ".jmp")); err != nil {
        t.Fatalf("expected config file to be created: %v", err)
    }
}

func contains(text, fragment string) bool {
    return len(text) >= len(fragment) && (text == fragment || contains(text[1:], fragment) || text[:len(fragment)] == fragment)
}
