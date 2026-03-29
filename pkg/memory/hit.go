package memory

import (
    "fmt"

    "github.com/nige123/go-jmp.nigelhamilton.com/pkg/file"
)

type Hit struct {
    file.Hit
}

func (h *Hit) Render() string {
    if h.MatchingText != "" {
        return fmt.Sprintf("    %s (%d) %s", h.RelativePath, h.LineNumber, h.MatchingText)
    }
    return fmt.Sprintf("    %s (%d)", h.RelativePath, h.LineNumber)
}
