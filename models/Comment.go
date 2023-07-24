package models

type Comment struct {
	CommentID    string  `gorm:"primaryKey" json:"comment_id"`
	Message      string  `json:"message"`
	CreatedUTC   float64 `json:"created_utc"`
	Score        uint8   `json:"score"`
	SubmissionID string  `json:"submission_id"`
}
