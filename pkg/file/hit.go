package file

import (
    "fmt"
    "os"
    "strings"
)

type Hit struct {
    MatchingText string `json:"matching-text,omitempty"`
    RelativePath string `json:"relative-path,omitempty"`
    AbsolutePath string `json:"absolute-path,omitempty"`
    LineNumber   int    `json:"line-number,omitempty"`
}

func NewHit(relativePath, absolutePath string, lineNumber int, matchingText string) *Hit {
    if lineNumber < 1 {
        lineNumber = 1
    }
    return &Hit{
        MatchingText: strings.ReplaceAll(matchingText, "\t", "    "),
        RelativePath: relativePath,
        AbsolutePath: absolutePath,
        LineNumber:   lineNumber,
    }
}

func (h *Hit) Render() string {
    if h.MatchingText != "" {
        return fmt.Sprintf("    (%d) %s", h.LineNumber, h.MatchingText)
    }
    return h.RelativePath
}

func (h *Hit) FileExists() bool {
    if h == nil || h.AbsolutePath == "" {
        return false
    }
    info, err := os.Stat(h.AbsolutePath)
    if err != nil {
        return false
    }
    return info.Mode().IsRegular()
}
