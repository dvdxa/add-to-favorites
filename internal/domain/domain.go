package domain

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Terminal struct {
	ID     int    `json:"id,omitempty"`
	Name   string `json:"name"`
	Status string `json:"status"`
}
type FakeTerminal struct {
	ID         int    `json:"id,omitempty"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	IsFavorite bool   `json:"is_favorite"`
}
