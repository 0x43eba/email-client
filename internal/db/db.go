package db

import (
    "database/sql"
    "time"

    _ "github.com/mattn/go-sqlite3"

    "demoproject.com/internal/model"
)

func LoadEmailsFromDB(dbFile string) (map[string][]model.Email, error) {
    db, err := sql.Open("sqlite3", dbFile)
    if err != nil {
        return nil, err
    }
    defer db.Close()

    folders := []string{"Inbox", "Sent", "Archive", "Deleted", "Spam"}
    emails := make(map[string][]model.Email)
    for _, folder := range folders {
        emails[folder] = []model.Email{}
    }

    rows, err := db.Query(`SELECT id, folder, sender, subject, date, body, "to" FROM emails ORDER BY id`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var (
            id      int
            folder  string
            sender  string
            subject string
            date    string
            body    string
            to      string
        )
        if err := rows.Scan(&id, &folder, &sender, &subject, &date, &body, &to); err != nil {
            return nil, err
        }
        e := model.Email{
            ID:      id,
            From:    sender,
            Subject: subject,
            Date:    date,
            Body:    body,
            To:      to,
        }
        emails[folder] = append(emails[folder], e)
    }

    for name, emSlice := range emails {
        for i, j := 0, len(emSlice)-1; i < j; i, j = i+1, j-1 {
            emSlice[i], emSlice[j] = emSlice[j], emSlice[i]
        }
        emails[name] = emSlice
    }

    return emails, nil
}

func UpdateEmailFolderInDB(emailID int, newFolder string) error {
    db, err := sql.Open("sqlite3", "emails.db")
    if err != nil {
        return err
    }
    defer db.Close()

    _, err = db.Exec("UPDATE emails SET folder = ? WHERE id = ?", newFolder, emailID)
    return err
}

func DeleteEmailFromDB(emailID int) error {
    db, err := sql.Open("sqlite3", "emails.db")
    if err != nil {
        return err
    }
    defer db.Close()

    _, err = db.Exec("DELETE FROM emails WHERE id = ?", emailID)
    return err
}

func InsertEmailIntoDB(to, subject, body string) error {
    db, err := sql.Open("sqlite3", "emails.db")
    if err != nil {
        return err
    }
    defer db.Close()

    now := time.Now().Format("2006-01-02 15:04")
    _, err = db.Exec(
        `INSERT INTO emails (folder, sender, subject, date, body, "to") VALUES (?, ?, ?, ?, ?, ?)`,
        "Sent", "Me", subject, now, body, to,
    )
    return err
}
