package template

import "testing"

func TestRenderSubstitutesToken(t *testing.T) {
    result := Render("hello [-name-]", map[string]string{"name": "world"})
    if result != "hello world" {
        t.Fatalf("expected %q, got %q", "hello world", result)
    }
}

func TestRenderSubstitutesMultipleTokens(t *testing.T) {
    result := Render("[-a-] and [-b-]", map[string]string{"a": "x", "b": "y"})
    if result != "x and y" {
        t.Fatalf("expected %q, got %q", "x and y", result)
    }
}

func TestRenderLeavesUnknownTokensUnchanged(t *testing.T) {
    result := Render("[-unknown-]", map[string]string{"other": "value"})
    if result != "[-unknown-]" {
        t.Fatalf("expected token to remain, got %q", result)
    }
}

func TestRenderNoTokensPassthrough(t *testing.T) {
    result := Render("no tokens here", map[string]string{"a": "b"})
    if result != "no tokens here" {
        t.Fatalf("expected %q, got %q", "no tokens here", result)
    }
}

func TestRenderEmptyInput(t *testing.T) {
    result := Render("", map[string]string{"a": "b"})
    if result != "" {
        t.Fatalf("expected empty string, got %q", result)
    }
}

func TestRenderEmptyParams(t *testing.T) {
    result := Render("[-token-]", map[string]string{})
    if result != "[-token-]" {
        t.Fatalf("expected token to remain, got %q", result)
    }
}

func TestRenderEditorCommandTemplate(t *testing.T) {
    tests := []struct {
        name     string
        template string
        params   map[string]string
        expected string
    }{
        {
            name:     "nano",
            template: `nano +[-line-number-] "[-filename-]"`,
            params:   map[string]string{"filename": "/tmp/test.go", "line-number": "42"},
            expected: `nano +42 "/tmp/test.go"`,
        },
        {
            name:     "vscode",
            template: `code -g [-filename-]:[-line-number-] &`,
            params:   map[string]string{"filename": "/tmp/test.go", "line-number": "10"},
            expected: `code -g /tmp/test.go:10 &`,
        },
        {
            name:     "vim",
            template: `vim +[-line-number-] [-filename-]`,
            params:   map[string]string{"filename": "/tmp/test.go", "line-number": "1"},
            expected: `vim +1 /tmp/test.go`,
        },
        {
            name:     "emacs",
            template: `emacs +[-line-number-] [-filename-]`,
            params:   map[string]string{"filename": "/home/user/main.go", "line-number": "99"},
            expected: `emacs +99 /home/user/main.go`,
        },
        {
            name:     "sublime",
            template: `subl [-filename-]:[-line-number-] &`,
            params:   map[string]string{"filename": "/tmp/app.py", "line-number": "5"},
            expected: `subl /tmp/app.py:5 &`,
        },
        {
            name:     "helix",
            template: `hx [-filename-]:[-line-number-]`,
            params:   map[string]string{"filename": "/tmp/main.rs", "line-number": "77"},
            expected: `hx /tmp/main.rs:77`,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Render(tt.template, tt.params)
            if result != tt.expected {
                t.Fatalf("expected %q, got %q", tt.expected, result)
            }
        })
    }
}

func TestRenderFindCommandTemplate(t *testing.T) {
    tests := []struct {
        name     string
        template string
        terms    string
        expected string
    }{
        {
            name:     "ripgrep",
            template: `rg --line-number --with-filename --no-heading --color never '[-search-terms-]'`,
            terms:    "TODO",
            expected: `rg --line-number --with-filename --no-heading --color never 'TODO'`,
        },
        {
            name:     "ag",
            template: `ag --nogroup '[-search-terms-]'`,
            terms:    "func main",
            expected: `ag --nogroup 'func main'`,
        },
        {
            name:     "git-grep",
            template: `git grep --full-name --untracked --text --line-number -e '[-search-terms-]'`,
            terms:    "import",
            expected: `git grep --full-name --untracked --text --line-number -e 'import'`,
        },
        {
            name:     "grep",
            template: `grep -rn '[-search-terms-]' .`,
            terms:    "error",
            expected: `grep -rn 'error' .`,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Render(tt.template, map[string]string{"search-terms": tt.terms})
            if result != tt.expected {
                t.Fatalf("expected %q, got %q", tt.expected, result)
            }
        })
    }
}

func TestRenderLocateCommandTemplate(t *testing.T) {
    tests := []struct {
        name     string
        template string
        terms    string
        expected string
    }{
        {
            name:     "locate",
            template: `locate '[-search-terms-]'`,
            terms:    "config.go",
            expected: `locate 'config.go'`,
        },
        {
            name:     "fd",
            template: `fd --type f '[-search-terms-]'`,
            terms:    "main",
            expected: `fd --type f 'main'`,
        },
        {
            name:     "find",
            template: `find / -name '*[-search-terms-]*' -type f 2>/dev/null`,
            terms:    "test",
            expected: `find / -name '*test*' -type f 2>/dev/null`,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Render(tt.template, map[string]string{"search-terms": tt.terms})
            if result != tt.expected {
                t.Fatalf("expected %q, got %q", tt.expected, result)
            }
        })
    }
}

func TestRenderBrowserCommandTemplate(t *testing.T) {
    tests := []struct {
        name     string
        template string
        url      string
        expected string
    }{
        {
            name:     "xdg-open",
            template: `xdg-open '[-url-]'`,
            url:      "https://example.com",
            expected: `xdg-open 'https://example.com'`,
        },
        {
            name:     "open-macos",
            template: `open '[-url-]'`,
            url:      "https://github.com/user/repo",
            expected: `open 'https://github.com/user/repo'`,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Render(tt.template, map[string]string{"url": tt.url})
            if result != tt.expected {
                t.Fatalf("expected %q, got %q", tt.expected, result)
            }
        })
    }
}

func TestRenderFilenameWithSpaces(t *testing.T) {
    result := Render(`nano +[-line-number-] "[-filename-]"`, map[string]string{
        "filename":    "/tmp/my file.go",
        "line-number": "1",
    })
    expected := `nano +1 "/tmp/my file.go"`
    if result != expected {
        t.Fatalf("expected %q, got %q", expected, result)
    }
}

func TestRenderSearchTermsWithSpaces(t *testing.T) {
    result := Render(`rg --line-number '[-search-terms-]'`, map[string]string{
        "search-terms": "func main",
    })
    expected := `rg --line-number 'func main'`
    if result != expected {
        t.Fatalf("expected %q, got %q", expected, result)
    }
}

func TestRenderSameTokenRepeated(t *testing.T) {
    result := Render("[-x-] [-x-]", map[string]string{"x": "hello"})
    if result != "hello hello" {
        t.Fatalf("expected %q, got %q", "hello hello", result)
    }
}

func TestRenderSubstitutionWithEmptyValue(t *testing.T) {
    result := Render("before [-token-] after", map[string]string{"token": ""})
    if result != "before  after" {
        t.Fatalf("expected %q, got %q", "before  after", result)
    }
}

func TestRenderNilParams(t *testing.T) {
    result := Render("[-token-]", nil)
    if result != "[-token-]" {
        t.Fatalf("expected token to remain, got %q", result)
    }
}
