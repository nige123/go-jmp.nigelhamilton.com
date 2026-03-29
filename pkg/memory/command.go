package memory

import (
    "fmt"
    "os"
    "os/exec"
)

type Command struct {
    CurrentDirectory string `json:"current-directory"`
    JmpCommand       string `json:"jmp-command"`
}

func (c *Command) Render() string {
    return fmt.Sprintf("%s> %s", c.CurrentDirectory, c.JmpCommand)
}

func (c *Command) Execute() error {
    if err := os.Chdir(c.CurrentDirectory); err != nil {
        return err
    }
    cmd := exec.Command("/bin/sh", "-c", c.JmpCommand)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}
