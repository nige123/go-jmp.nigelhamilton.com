package finder

import (
    "testing"

    "github.com/nige123/go-jmp.nigelhamilton.com/pkg/file"
)

type fakeConfig struct{}

func (fakeConfig) GetTemplate(key string, params map[string]string) (string, error) {
    return "fake", nil
}

func TestFindFilesOnFilesystem(t *testing.T) {
    finder := NewWithRunners(fakeConfig{}, func(command string) ([]string, error) {
        return []string{
            "/usr/share/doc/README.md",
            "/home/user/projects/README.md",
        }, nil
    }, nil)

    hits, err := finder.FindFilesOnFilesystem("README.md")
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    if len(hits) != 2 {
        t.Fatalf("expected 2 hits, got %d", len(hits))
    }
}

func TestFindInFilesIgnoresMalformedLines(t *testing.T) {
    finder := NewWithRunners(fakeConfig{}, func(command string) ([]string, error) {
        return []string{
            "testdata/1.txt",
            "1:wrong-heading-style-match",
            "testdata/1.txt:1:the real match",
        }, nil
    }, nil)

    hits, err := finder.FindInFiles("ignored")
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    if len(hits) != 2 {
        t.Fatalf("expected 2 hits, got %d", len(hits))
    }

    first, ok := hits[0].(*file.Hit)
    if !ok {
        t.Fatalf("expected first hit type *file.Hit")
    }
    second, ok := hits[1].(*file.Hit)
    if !ok {
        t.Fatalf("expected second hit type *file.Hit")
    }

    if first.RelativePath != "testdata/1.txt" || first.LineNumber != 1 || first.MatchingText != "" {
        t.Fatalf("unexpected first hit: %#v", first)
    }
    if second.RelativePath != "testdata/1.txt" || second.LineNumber != 1 || second.MatchingText != "the real match" {
        t.Fatalf("unexpected second hit: %#v", second)
    }
}
