package file

import (
    "os"
    "path/filepath"
    "regexp"
    "strconv"
    "strings"
    "unicode"
)

type HitLater struct {
    Hit
    TextToMatch string
}

var (
    perlLinePattern = regexp.MustCompile(`at\s+(\S+)\s+line\s+(\d+)`)
    rakuLinePattern = regexp.MustCompile(`at\s+([^\s:]+):(\d+)`)
    rakuColonModule = regexp.MustCompile(`at\s+([^\s:]+)\s*\(\S+\):(\d+)`)
    rakuLineModule  = regexp.MustCompile(`at\s+([^\s:]+)\s*\(\S+\)\s+line\s+(\d+)`)
    ansiPattern     = regexp.MustCompile(`\[[0-9;]*m`)
)

func NewHitLater(text string) *HitLater {
    cleaned := stripControlCodes(text)
    cleaned = ansiPattern.ReplaceAllString(cleaned, "")
    return &HitLater{TextToMatch: cleaned}
}

func (h *HitLater) Render() string {
    return h.TextToMatch
}

func (h *HitLater) FindFilePath() {
    if h.FileExists() {
        return
    }

    if match := perlLinePattern.FindStringSubmatch(h.TextToMatch); len(match) == 3 {
        if h.foundFilePath(match[1], parseLine(match[2])) {
            return
        }
    }

    if match := rakuLinePattern.FindStringSubmatch(h.TextToMatch); len(match) == 3 {
        if h.foundFilePath(match[1], parseLine(match[2])) {
            return
        }
    }

    if match := rakuColonModule.FindStringSubmatch(h.TextToMatch); len(match) == 3 {
        if h.foundFilePath(match[1], parseLine(match[2])) {
            return
        }
    }

    if match := rakuLineModule.FindStringSubmatch(h.TextToMatch); len(match) == 3 {
        if h.foundFilePath(match[1], parseLine(match[2])) {
            return
        }
    }

    for _, token := range strings.Fields(h.TextToMatch) {
        if h.foundFilePath(token, 1) {
            return
        }
    }
}

func (h *HitLater) foundFilePath(path string, lineNumber int) bool {
    info, err := os.Stat(path)
    if err != nil || !info.Mode().IsRegular() {
        return false
    }
    abs, err := filepath.Abs(path)
    if err != nil {
        return false
    }
    h.RelativePath = path
    h.AbsolutePath = abs
    h.LineNumber = lineNumber
    h.MatchingText = h.TextToMatch
    return true
}

func parseLine(raw string) int {
    value, err := strconv.Atoi(raw)
    if err != nil || value < 1 {
        return 1
    }
    return value
}

func stripControlCodes(input string) string {
    var b strings.Builder
    b.Grow(len(input))
    for _, r := range input {
        if unicode.IsControl(r) && r != '\n' && r != '\r' && r != '\t' {
            continue
        }
        b.WriteRune(r)
    }
    return b.String()
}
