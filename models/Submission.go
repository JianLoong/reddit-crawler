package models

type Submission struct {
	// gorm.Model
	SubmissionID string    `gorm:"primaryKey" json:"submission_id"`
	Title        string    `json:"title"`
	Selftext     string    `json:"selftext"`
	CreatedUTC   float64   `json:"created_utc"`
	Permalink    string    `json:"permalink"`
	Score        uint8     `json:"score"`
	Url          string    `json:"url"`
	Comments     []Comment `json:"comments"`
}
