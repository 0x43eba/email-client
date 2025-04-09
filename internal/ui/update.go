package ui

import (
    "fmt"

    tea "github.com/charmbracelet/bubbletea"

    "demoproject.com/internal/db"
    "demoproject.com/internal/model"
    "unicode/utf8"
)

type refreshMsg struct{}

func RefreshCmd() tea.Cmd {
    return func() tea.Msg {
        return refreshMsg{}
    }
}

func NewProgramModel(m model.Model) *programModel {
    return &programModel{m}
}

type programModel struct {
    model model.Model
}

func (p *programModel) Init() tea.Cmd {
    return nil
}

func (p *programModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    m := &p.model
    if m.Mode == model.ModeCompose {
        switch keyMsg := msg.(type) {
        case tea.KeyMsg:
            switch keyMsg.String() {
            case "esc":
                m.Mode = model.ModeList

            case "tab":
                m.ComposeField = (m.ComposeField + 1) % 3

            case "shift+tab":
                m.ComposeField = (m.ComposeField + 2) % 3

            case "enter":
                if m.ComposeField == model.FieldBody {
                    m.ComposeBody += "\n"
                } else {
                    m.ComposeField = (m.ComposeField + 1) % 3
                }

            case "backspace":
                switch m.ComposeField {
                case model.FieldTo:
                    if len(m.ComposeTo) > 0 {
                        _, size := utf8.DecodeLastRuneInString(m.ComposeTo)
                        m.ComposeTo = m.ComposeTo[:len(m.ComposeTo)-size]
                    }
                case model.FieldSubject:
                    if len(m.ComposeSubject) > 0 {
                        _, size := utf8.DecodeLastRuneInString(m.ComposeSubject)
                        m.ComposeSubject = m.ComposeSubject[:len(m.ComposeSubject)-size]
                    }
                case model.FieldBody:
                    if len(m.ComposeBody) > 0 {
                        _, size := utf8.DecodeLastRuneInString(m.ComposeBody)
                        m.ComposeBody = m.ComposeBody[:len(m.ComposeBody)-size]
                    }
                }

            case "ctrl+s":
                if err := db.InsertEmailIntoDB(m.ComposeTo, m.ComposeSubject, m.ComposeBody); err != nil {
                    fmt.Printf("Error sending email: %v\n", err)
                } else {
                    sentEmails, err := db.LoadEmailsFromDB("emails.db")
                    if err == nil {
                        m.Emails["Sent"] = sentEmails["Sent"]
                    } else {
                        fmt.Printf("Error reloading sent emails: %v\n", err)
                    }
                }
                for i, folder := range m.Folders {
                    if folder == "Sent" {
                        m.FolderCursor = i
                        m.EmailCursor = 0
                        break
                    }
                }
                m.Mode = model.ModeList

            default:
                switch m.ComposeField {
                case model.FieldTo:
                    m.ComposeTo += keyMsg.String()
                case model.FieldSubject:
                    m.ComposeSubject += keyMsg.String()
                case model.FieldBody:
                    m.ComposeBody += keyMsg.String()
                }
            }
        }
        return p, RefreshCmd()
    }

    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.Width = msg.Width

    case tea.KeyMsg:
        switch msg.String() {
        case "q":
            return p, tea.Quit

        case "s":
            m.Mode = model.ModeCompose
            m.ComposeTo, m.ComposeSubject, m.ComposeBody = "", "", ""
            m.ComposeField = model.FieldTo
            return p, RefreshCmd()

        case "r", "R":
            newEmails, err := db.LoadEmailsFromDB("emails.db")
            if err == nil {
                m.Emails = newEmails
            } else {
                fmt.Printf("Error refreshing emails: %v\n", err)
            }
            return p, RefreshCmd()

        case "A", "a":
            currentFolder := m.Folders[m.FolderCursor]
            emailList := m.Emails[currentFolder]
            if len(emailList) == 0 {
                break
            }
            selectedEmail := emailList[m.EmailCursor]
            currentAction := m.Actions[m.ActionCursor]

            var err error
            switch currentAction {
            case "Reply":
                m.Mode = model.ModeCompose
                m.ComposeTo = selectedEmail.From
                m.ComposeSubject = "Re: " + selectedEmail.Subject
                m.ComposeBody = "\n--- Original Message ---\n" + selectedEmail.Body
                m.ComposeField = model.FieldBody
                return p, RefreshCmd()

            case "Delete":
                if currentFolder != "Deleted" {
                    err = db.UpdateEmailFolderInDB(selectedEmail.ID, "Deleted")
                    if err == nil {
                        m.Emails["Deleted"] = append(m.Emails["Deleted"], selectedEmail)
                    }
                } else {
                    err = db.DeleteEmailFromDB(selectedEmail.ID)
                }
                m.Emails[currentFolder] = model.RemoveEmail(m.Emails[currentFolder], m.EmailCursor)
                if m.EmailCursor >= len(m.Emails[currentFolder]) && m.EmailCursor > 0 {
                    m.EmailCursor--
                }

            case "Archive":
                if currentFolder == "Inbox" {
                    err = db.UpdateEmailFolderInDB(selectedEmail.ID, "Archive")
                    if err == nil {
                        m.Emails["Archive"] = append(m.Emails["Archive"], selectedEmail)
                        m.Emails[currentFolder] = model.RemoveEmail(m.Emails[currentFolder], m.EmailCursor)
                    }
                }
                if m.EmailCursor >= len(m.Emails[currentFolder]) && m.EmailCursor > 0 {
                    m.EmailCursor--
                }

            case "Move to Inbox":
                if currentFolder == "Archive" {
                    err = db.UpdateEmailFolderInDB(selectedEmail.ID, "Inbox")
                    if err == nil {
                        m.Emails["Inbox"] = append(m.Emails["Inbox"], selectedEmail)
                        m.Emails[currentFolder] = model.RemoveEmail(m.Emails[currentFolder], m.EmailCursor)
                    }
                }
                if m.EmailCursor >= len(m.Emails[currentFolder]) && m.EmailCursor > 0 {
                    m.EmailCursor--
                }

            case "Restore":
                if currentFolder == "Deleted" {
                    err = db.UpdateEmailFolderInDB(selectedEmail.ID, "Inbox")
                    if err == nil {
                        m.Emails["Inbox"] = append(m.Emails["Inbox"], selectedEmail)
                        m.Emails[currentFolder] = model.RemoveEmail(m.Emails[currentFolder], m.EmailCursor)
                    }
                }
                if m.EmailCursor >= len(m.Emails[currentFolder]) && m.EmailCursor > 0 {
                    m.EmailCursor--
                }

            case "Delete Permanently":
                if currentFolder == "Deleted" {
                    err = db.DeleteEmailFromDB(selectedEmail.ID)
                    if err == nil {
                        m.Emails[currentFolder] = model.RemoveEmail(m.Emails[currentFolder], m.EmailCursor)
                    }
                }
                if m.EmailCursor >= len(m.Emails[currentFolder]) && m.EmailCursor > 0 {
                    m.EmailCursor--
                }

            case "Mark as Not Spam":
                if currentFolder == "Spam" {
                    err = db.UpdateEmailFolderInDB(selectedEmail.ID, "Inbox")
                    if err == nil {
                        m.Emails["Inbox"] = append(m.Emails["Inbox"], selectedEmail)
                        m.Emails[currentFolder] = model.RemoveEmail(m.Emails[currentFolder], m.EmailCursor)
                    }
                }
                if m.EmailCursor >= len(m.Emails[currentFolder]) && m.EmailCursor > 0 {
                    m.EmailCursor--
                }

            case "Mark as Spam":
                err = db.UpdateEmailFolderInDB(selectedEmail.ID, "Spam")
                if err == nil {
                    m.Emails["Spam"] = append(m.Emails["Spam"], selectedEmail)
                    m.Emails[currentFolder] = model.RemoveEmail(m.Emails[currentFolder], m.EmailCursor)
                }
                if m.EmailCursor >= len(m.Emails[currentFolder]) && m.EmailCursor > 0 {
                    m.EmailCursor--
                }

            default:
                // For "Open", "Mark Read", etc., no special DB action here.
            }

            if err != nil {
                fmt.Printf("Database update error: %v\n", err)
            }
            return p, RefreshCmd()

        case "left":
            if m.FolderCursor > 0 {
                m.FolderCursor--
                m.EmailCursor = 0
                currentFolder := m.Folders[m.FolderCursor]
                m.Actions = model.GetActionsForFolder(currentFolder)
                m.ActionCursor = 0
            }

        case "right":
            if m.FolderCursor < len(m.Folders)-1 {
                m.FolderCursor++
                m.EmailCursor = 0
                currentFolder := m.Folders[m.FolderCursor]
                m.Actions = model.GetActionsForFolder(currentFolder)
                m.ActionCursor = 0
            }

        case "up":
            if m.EmailCursor > 0 {
                m.EmailCursor--
            }

        case "down":
            currentFolder := m.Folders[m.FolderCursor]
            if m.EmailCursor < len(m.Emails[currentFolder])-1 {
                m.EmailCursor++
            }

        case "z":
            if m.ActionCursor > 0 {
                m.ActionCursor--
            }

        case "x":
            if m.ActionCursor < len(m.Actions)-1 {
                m.ActionCursor++
            }
        }

    case refreshMsg:
        // Force re-render.

    }
    return p, RefreshCmd()
}
