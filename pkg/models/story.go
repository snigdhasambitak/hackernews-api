package models

type Story struct {
	Author   string `json:"author"`
	Karma    int    `json:"karma"`
	Comments int    `json:"comments"`
	Title    string `json:"title"`
	Position int    `json:"position"`
}
