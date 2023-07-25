package models

type Index struct {
	CreatedUTC   float64 `json:"created_utc"`
	Score        uint8   `json:"score"`
	SubmissionID string  `json:"submission_id"`
}
