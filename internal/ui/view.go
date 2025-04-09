package ui

import (
    "fmt"

    "demoproject.com/internal/model"
)

func (p *programModel) View() string {
    m := p.model

    if m.Mode == model.ModeCompose {
        return composeView(m)
    }
    return listView(m)
}

func composeView(m model.Model) string {
    var s string
    s += "Compose Email\n"
    s += "-------------\n"

    fieldIndicator := func(idx int, label, value string) string {
        marker := "  "
        if m.ComposeField == idx {
            marker = "> "
        }
        return fmt.Sprintf("%s%s: %s\n", marker, label, value)
    }

    s += fieldIndicator(model.FieldTo, "To", m.ComposeTo)
    s += fieldIndicator(model.FieldSubject, "Subject", m.ComposeSubject)
    s += fieldIndicator(model.FieldBody, "Body", m.ComposeBody)

    s += "\nPress Tab to cycle fields, Ctrl+S (in Body) to send, Enter in Body to add a newline, or Esc to cancel."

    return model.CenterText(s, m.Width)
}

func listView(m model.Model) string {
    var s string

    s += "Folders: "
    for i, folder := range m.Folders {
        if i == m.FolderCursor {
            s += fmt.Sprintf("[ %s ] ", folder)
        } else {
            s += fmt.Sprintf("  %s   ", folder)
        }
    }
    s += "\n"

    s += "Actions: "
    for i, action := range m.Actions {
        if i == m.ActionCursor {
            s += fmt.Sprintf("< %s >", action)
        } else {
            s += fmt.Sprintf("[ %s ]", action)
        }
        if i < len(m.Actions)-1 {
            s += "  "
        }
    }

    s += "\n\nCommands: ←/→: Change folders | ↑/↓: Navigate emails | z/x: Select action | A: Execute action | s: Compose email | r: Refresh | q: Quit\n\n"

    currentFolder := m.Folders[m.FolderCursor]
    emailList := m.Emails[currentFolder]
    s += fmt.Sprintf("== %s ==\n", currentFolder)
    s += "------------------------------------------------------\n"

    if len(emailList) == 0 {
        s += "No emails in this folder.\n"
    } else {
        s += fmt.Sprintf("%-15s  %-25s  %s\n", "From", "Subject", "Date")
        s += "------------------------------------------------------\n"
        for i, email := range emailList {
            cursor := " "
            if m.EmailCursor == i {
                cursor = ">"
            }
            s += fmt.Sprintf("%s %-14s  %-25s  %s\n",
                cursor, email.From, email.Subject, email.Date)
        }
    }

    s += "\n------------------------------------------------------\n"
    s += "Reading Pane\n"
    s += "------------------------------------------------------\n"

    if len(emailList) > 0 {
        s += emailList[m.EmailCursor].Body
    } else {
        s += "Select another folder or add emails to view details."
    }

    return s
}
