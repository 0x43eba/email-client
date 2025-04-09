package model

func RemoveEmail(emails []Email, i int) []Email {
    return append(emails[:i], emails[i+1:]...)
}
