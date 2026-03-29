package main

import (
    "fmt"
    "os"

    "github.com/nige123/go-jmp.nigelhamilton.com/internal/app"
)

func main() {
    application, err := app.New(os.Args[1:])
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    if err := application.Run(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
