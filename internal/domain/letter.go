package domain

type Letter struct {
	ID          int    `json:"id"`
	Content     string `json:"content"`
	OpenedCount int    `json:"opened_count"`
	LastOpened  string `json:"last_opened"`
	Date        string `json:"show_date"`
	LikeCount   int    `json:"like_count"`
}
