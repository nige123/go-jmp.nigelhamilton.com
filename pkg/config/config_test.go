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

func TestGetTemplateEditorCommand(t *testing.T) {
    home := t.TempDir()
    t.Setenv("HOME", home)

    cfg, err := New()
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }

    rendered, err := cfg.GetTemplate("editor.command.template", map[string]string{
        "filename":    "/tmp/test.go",
        "line-number": "42",
    })
    if err != nil {
        t.Fatalf("expected rendered template, got %v", err)
    }

    // Default template no longer has quotes — editor.go adds them around filename
    expected := `nano +42 /tmp/test.go`
    if rendered != expected {
        t.Fatalf("expected %q, got %q", expected, rendered)
    }
}

func TestGetTemplateLocateCommand(t *testing.T) {
    home := t.TempDir()
    t.Setenv("HOME", home)

    cfg, err := New()
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }

    rendered, err := cfg.GetTemplate("locate.command.template", map[string]string{
        "search-terms": "config.go",
    })
    if err != nil {
        t.Fatalf("expected rendered template, got %v", err)
    }

    expected := "locate 'config.go'"
    if rendered != expected {
        t.Fatalf("expected %q, got %q", expected, rendered)
    }
}

func TestGetTemplateBrowserCommand(t *testing.T) {
    home := t.TempDir()
    t.Setenv("HOME", home)

    cfg, err := New()
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }

    rendered, err := cfg.GetTemplate("browser.command.template", map[string]string{
        "url": "https://example.com",
    })
    if err != nil {
        t.Fatalf("expected rendered template, got %v", err)
    }

    expected := "xdg-open 'https://example.com'"
    if rendered != expected {
        t.Fatalf("expected %q, got %q", expected, rendered)
    }
}

func TestGetTemplateCustomConfig(t *testing.T) {
    home := t.TempDir()
    t.Setenv("HOME", home)

    customConfig := `
editor.command.template = code -g [-filename-]:[-line-number-] &
find.command.template = ag --nogroup '[-search-terms-]'
`
    if err := os.WriteFile(filepath.Join(home, ".jmp"), []byte(customConfig), 0o644); err != nil {
        t.Fatalf("could not write custom config: %v", err)
    }

    cfg, err := New()
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }

    editorCmd, err := cfg.GetTemplate("editor.command.template", map[string]string{
        "filename":    "/home/user/main.go",
        "line-number": "10",
    })
    if err != nil {
        t.Fatalf("expected rendered template, got %v", err)
    }
    expectedEditor := "code -g /home/user/main.go:10 &"
    if editorCmd != expectedEditor {
        t.Fatalf("expected %q, got %q", expectedEditor, editorCmd)
    }

    findCmd, err := cfg.GetTemplate("find.command.template", map[string]string{
        "search-terms": "func main",
    })
    if err != nil {
        t.Fatalf("expected rendered template, got %v", err)
    }
    expectedFind := "ag --nogroup 'func main'"
    if findCmd != expectedFind {
        t.Fatalf("expected %q, got %q", expectedFind, findCmd)
    }
}

func TestDoubleQuotesStrippedFromConfigValues(t *testing.T) {
    home := t.TempDir()
    t.Setenv("HOME", home)

    // Double quotes in config values should be stripped — we handle quoting ourselves
    configContent := `editor.command.template = "nano +[-line-number-] [-filename-]"
find.command.template = "rg --line-number '[-search-terms-]'"
locate.command.template = "locate '[-search-terms-]'"
`
    if err := os.WriteFile(filepath.Join(home, ".jmp"), []byte(configContent), 0o644); err != nil {
        t.Fatalf("could not write config: %v", err)
    }

    cfg, err := New()
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }

    editorCmd, err := cfg.GetTemplate("editor.command.template", map[string]string{
        "filename":    "/tmp/test.go",
        "line-number": "42",
    })
    if err != nil {
        t.Fatalf("expected rendered template, got %v", err)
    }
    if editorCmd != "nano +42 /tmp/test.go" {
        t.Fatalf("expected %q, got %q", "nano +42 /tmp/test.go", editorCmd)
    }

    findCmd, err := cfg.GetTemplate("find.command.template", map[string]string{
        "search-terms": "TODO",
    })
    if err != nil {
        t.Fatalf("expected rendered template, got %v", err)
    }
    if findCmd != "rg --line-number 'TODO'" {
        t.Fatalf("expected %q, got %q", "rg --line-number 'TODO'", findCmd)
    }

    locateCmd, err := cfg.GetTemplate("locate.command.template", map[string]string{
        "search-terms": "config.go",
    })
    if err != nil {
        t.Fatalf("expected rendered template, got %v", err)
    }
    if locateCmd != "locate 'config.go'" {
        t.Fatalf("expected %q, got %q", "locate 'config.go'", locateCmd)
    }
}

func TestGetTemplateMissingKey(t *testing.T) {
    home := t.TempDir()
    t.Setenv("HOME", home)

    cfg, err := New()
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }

    _, err = cfg.GetTemplate("nonexistent.template", map[string]string{"a": "b"})
    if err == nil {
        t.Fatalf("expected error for missing key, got nil")
    }
}

func contains(text, fragment string) bool {
    return len(text) >= len(fragment) && (text == fragment || contains(text[1:], fragment) || text[:len(fragment)] == fragment)
}
