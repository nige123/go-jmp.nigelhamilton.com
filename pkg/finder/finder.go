package finder

import (
    "bytes"
    "errors"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "regexp"
    "strconv"
    "strings"

    "github.com/nige123/go-jmp.nigelhamilton.com/pkg/file"
    "github.com/nige123/go-jmp.nigelhamilton.com/pkg/model"
)

type ConfigProvider interface {
    GetTemplate(key string, params map[string]string) (string, error)
}

type Finder struct {
    config      ConfigProvider
    runStdout   func(command string) ([]string, error)
    runCombined func(command string) ([]string, error)
}

var findOutputPattern = regexp.MustCompile(`^(.*):(\d+):(.*)$`)

func New(config ConfigProvider) *Finder {
    return &Finder{
        config:      config,
        runStdout:   runStdoutShell,
        runCombined: runCombinedShell,
    }
}

func NewWithRunners(config ConfigProvider, runStdout func(string) ([]string, error), runCombined func(string) ([]string, error)) *Finder {
    finder := New(config)
    if runStdout != nil {
        finder.runStdout = runStdout
    }
    if runCombined != nil {
        finder.runCombined = runCombined
    }
    return finder
}

func (f *Finder) FindLineInFile(filename string, lineNumber int) *file.Hit {
    abs, _ := filepath.Abs(filename)
    return file.NewHit(filename, abs, lineNumber, "")
}

func (f *Finder) FindMatchingLineInFile(filename, searchTerms string) (*file.Hit, error) {
    abs, _ := filepath.Abs(filename)
    hit := file.NewHit(filename, abs, 1, "")

    contents, err := os.ReadFile(filename)
    if err != nil {
        return hit, nil
    }

    for index, line := range strings.Split(strings.ReplaceAll(string(contents), "\r\n", "\n"), "\n") {
        if strings.Contains(line, searchTerms) {
            hit.LineNumber = index + 1
            hit.MatchingText = line
            return hit, nil
        }
    }
    return hit, nil
}

func (f *Finder) FindInFiles(searchTerms string) ([]model.Renderable, error) {
    findCommand, err := f.config.GetTemplate("find.command.template", map[string]string{"search-terms": searchTerms})
    if err != nil {
        return nil, err
    }

    lines, err := f.runStdout(findCommand)
    if err != nil {
        return nil, err
    }

    hits := make([]model.Renderable, 0, len(lines))
    previousFile := ""
    for _, line := range lines {
        filePath, lineNumber, matchingText, ok := parseFindOutputLine(line)
        if !ok {
            continue
        }

        abs, _ := filepath.Abs(filePath)
        if filePath != previousFile {
            hits = append(hits, file.NewHit(filePath, abs, 1, ""))
            previousFile = filePath
        }
        hits = append(hits, file.NewHit(filePath, abs, lineNumber, matchingText))
    }

    return hits, nil
}

func (f *Finder) FindFilesInCommandOutput(command string) ([]model.Renderable, error) {
    lines, err := f.runCombined(command)
    if err != nil {
        return nil, err
    }
    hits := make([]model.Renderable, 0, len(lines))
    for _, line := range lines {
        hits = append(hits, file.NewHitLater(line))
    }
    return hits, nil
}

func parseFindOutputLine(line string) (string, int, string, bool) {
    match := findOutputPattern.FindStringSubmatch(line)
    if len(match) != 4 {
        return "", 0, "", false
    }
    lineNumber, err := strconv.Atoi(match[2])
    if err != nil {
        return "", 0, "", false
    }
    return match[1], lineNumber, match[3], true
}

func runStdoutShell(command string) ([]string, error) {
    cmd := exec.Command("/bin/sh", "-c", command)
    var stdout bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr

    err := cmd.Run()
    if err != nil {
        var exitErr *exec.ExitError
        if !errors.As(err, &exitErr) {
            return nil, fmt.Errorf("could not run search command: %w", err)
        }
    }
    return splitLines(stdout.String()), nil
}

func runCombinedShell(command string) ([]string, error) {
    cmd := exec.Command("/bin/sh", "-c", command)
    var output bytes.Buffer
    cmd.Stdout = &output
    cmd.Stderr = &output

    err := cmd.Run()
    if err != nil {
        var exitErr *exec.ExitError
        if !errors.As(err, &exitErr) {
            return nil, fmt.Errorf("could not run command output search: %w", err)
        }
    }
    return splitLines(output.String()), nil
}

func splitLines(output string) []string {
    if output == "" {
        return nil
    }
    normalized := strings.ReplaceAll(output, "\r\n", "\n")
    normalized = strings.TrimSuffix(normalized, "\n")
    if normalized == "" {
        return nil
    }
    return strings.Split(normalized, "\n")
}
