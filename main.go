package main

import (
    "fmt"
    "os"

    tea "github.com/charmbracelet/bubbletea"

    "demoproject.com/internal/db"
    "demoproject.com/internal/model"
    "demoproject.com/internal/ui"
)

func main() {
    fmt.Print("\033[H\033[2J")

    emails, err := db.LoadEmailsFromDB("emails.db")
    if err != nil {
        fmt.Printf("Error loading emails from db: %v\n", err)
        os.Exit(1)
    }

    m := model.NewModel(emails)

    p := tea.NewProgram(ui.NewProgramModel(m))
    if _, err := p.Run(); err != nil {
        fmt.Printf("Error running program: %v\n", err)
        os.Exit(1)
    }

    fmt.Print("\033[H\033[2J")
}
