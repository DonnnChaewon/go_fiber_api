package main

type Status string

const (
	Watched Status = "watched"
	Watching Status = "watching"
	ToWatch Status = "to_watch"
)

type User struct {
	ID uint `json:"id" gorm:"primaryKey"`
	Username string `json:"username"`
	Password string `json:"-"`
}

type Movie struct {
	ID uint `json:"id" gorm:"primaryKey"`
	Title string `json:"title"`
	Status Status `json:"status" gorm:"default:to_watch"`
	Director string `json:"director"`
	UserID int `json:"user_id"`
}