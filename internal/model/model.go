package model

import "strings"

type modeType int

const (
    ModeList modeType = iota
    ModeCompose
)


const (
    FieldTo = iota
    FieldSubject
    FieldBody
)

type Model struct {

    Folders      []string
    FolderCursor int
    Actions      []string
    Emails       map[string][]Email
    EmailCursor  int
    ActionCursor int
    Width        int

    Mode           modeType
    ComposeTo      string
    ComposeSubject string
    ComposeBody    string
    ComposeField   int
}

func NewModel(emails map[string][]Email) Model {
    folders := []string{"Inbox", "Sent", "Archive", "Deleted", "Spam"}
    actions := GetActionsForFolder(folders[0])

    return Model{
        Folders:      folders,
        FolderCursor: 0,
        Actions:      actions,
        Emails:       emails,
        EmailCursor:  0,
        ActionCursor: 0,
        Width:        80,
        Mode:         ModeList,
    }
}

func GetActionsForFolder(folder string) []string {
    switch folder {
    case "Inbox":
        return []string{"Open", "Reply", "Mark Read", "Delete", "Archive", "Mark as Spam"}
    case "Sent":
        return []string{"Open", "Mark Read", "Delete", "Mark as Spam"}
    case "Archive":
        return []string{"Open", "Mark Read", "Delete", "Move to Inbox", "Mark as Spam"}
    case "Deleted":
        return []string{"Open", "Restore", "Delete Permanently", "Mark as Spam"}
    case "Spam":
        return []string{"Open", "Mark as Not Spam", "Delete"}
    default:
        return []string{"Open"}
    }
}

func CenterText(s string, width int) string {
    lines := strings.Split(s, "\n")
    var centered []string
    for _, line := range lines {
        pad := (width - len(line)) / 2
        if pad < 0 {
            pad = 0
        }
        centered = append(centered, linePadding(pad)+line)
    }
    return strings.Join(centered, "\n")
}

func linePadding(n int) string {
    if n <= 0 {
        return ""
    }
    return string(make([]byte, n))
}
