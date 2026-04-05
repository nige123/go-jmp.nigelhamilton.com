package editor

import (
    "fmt"
    "os"
    "os/exec"
    "strconv"

    "github.com/nige123/go-jmp.nigelhamilton.com/pkg/file"
)

type ConfigProvider interface {
    GetTemplate(key string, params map[string]string) (string, error)
}

type MemorySaver interface {
    Save(jmpCommand string, hit *file.Hit) error
}

type Editor struct {
    config ConfigProvider
    memory MemorySaver
    runner func(command string) error
}

func New(config ConfigProvider, memory MemorySaver) *Editor {
    return &Editor{config: config, memory: memory, runner: runForegroundCommand}
}

func NewWithRunner(config ConfigProvider, memory MemorySaver, runner func(command string) error) *Editor {
    e := New(config, memory)
    if runner != nil {
        e.runner = runner
    }
    return e
}

func (e *Editor) Edit(jmpCommand string, hit *file.Hit) error {
    if hit == nil || !hit.FileExists() {
        return nil
    }
    if err := e.EditAtLine(hit.AbsolutePath, hit.LineNumber); err != nil {
        return err
    }
    return e.memory.Save(jmpCommand, hit)
}

func (e *Editor) EditAtLine(filename string, lineNumber int) error {
    command, err := e.config.GetTemplate("editor.command.template", map[string]string{
        "filename":    `"` + filename + `"`,
        "line-number": strconv.Itoa(lineNumber),
    })
    if err != nil {
        return err
    }
    if err := e.runner(command); err != nil {
        return fmt.Errorf("could not execute editor command: %w", err)
    }
    return nil
}

func (e *Editor) EditFile(filename string) error {
    return e.EditAtLine(filename, 1)
}

func runForegroundCommand(command string) error {
    cmd := exec.Command("/bin/sh", "-c", command)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}
