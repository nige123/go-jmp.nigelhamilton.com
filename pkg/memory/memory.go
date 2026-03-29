package memory

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"

    "github.com/nige123/go-jmp.nigelhamilton.com/pkg/file"
    "github.com/nige123/go-jmp.nigelhamilton.com/pkg/model"
)

type Memory struct {
    maxEntries int
    file       string
    latestJmps []map[string]interface{}
}

func New() (*Memory, error) {
    home, err := os.UserHomeDir()
    if err != nil {
        return nil, fmt.Errorf("could not determine HOME directory: %w", err)
    }

    m := &Memory{maxEntries: 100, file: filepath.Join(home, ".jmp.hist")}
    if _, err := os.Stat(m.file); os.IsNotExist(err) {
        return m, nil
    }

    payload, err := os.ReadFile(m.file)
    if err != nil {
        return nil, fmt.Errorf("could not read memory file %s: %w", m.file, err)
    }
    if len(payload) == 0 {
        return m, nil
    }

    if err := json.Unmarshal(payload, &m.latestJmps); err != nil {
        return nil, fmt.Errorf("could not parse memory file %s: %w", m.file, err)
    }
    return m, nil
}

func (m *Memory) GetRecentJmps(lastNJmps int) []model.Renderable {
    if lastNJmps < 0 {
        lastNJmps = 0
    }
    if lastNJmps > len(m.latestJmps) {
        lastNJmps = len(m.latestJmps)
    }

    hits := make([]model.Renderable, 0, lastNJmps)
    for _, raw := range m.latestJmps[:lastNJmps] {
        if _, ok := raw["jmp-command"]; ok {
            hits = append(hits, &Command{
                CurrentDirectory: asString(raw["current-directory"]),
                JmpCommand:       asString(raw["jmp-command"]),
            })
            continue
        }
        hits = append(hits, &Hit{Hit: file.Hit{
            LineNumber:   asInt(raw["line-number"]),
            RelativePath: asString(raw["relative-path"]),
            AbsolutePath: asString(raw["absolute-path"]),
            MatchingText: asString(raw["matching-text"]),
        }})
    }
    return hits
}

func (m *Memory) Save(jmpCommand string, hit *file.Hit) error {
    cwd, err := os.Getwd()
    if err != nil {
        return fmt.Errorf("could not determine working directory: %w", err)
    }

    commandRecord := map[string]interface{}{
        "current-directory": cwd,
        "jmp-command":       jmpCommand,
    }
    hitRecord := map[string]interface{}{
        "line-number":   hit.LineNumber,
        "relative-path": hit.RelativePath,
        "absolute-path": hit.AbsolutePath,
        "matching-text": hit.MatchingText,
    }

    m.latestJmps = append([]map[string]interface{}{hitRecord, commandRecord}, m.latestJmps...)
    if len(m.latestJmps) > m.maxEntries {
        m.latestJmps = m.latestJmps[:m.maxEntries]
    }

    payload, err := json.Marshal(m.latestJmps)
    if err != nil {
        return fmt.Errorf("could not serialize memory entries: %w", err)
    }
    if err := os.WriteFile(m.file, payload, 0o644); err != nil {
        return fmt.Errorf("could not write memory file %s: %w", m.file, err)
    }

    return nil
}

func asString(value interface{}) string {
    if value == nil {
        return ""
    }
    if s, ok := value.(string); ok {
        return s
    }
    return ""
}

func asInt(value interface{}) int {
    switch v := value.(type) {
    case float64:
        return int(v)
    case int:
        return v
    default:
        return 1
    }
}
