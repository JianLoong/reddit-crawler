package models

type SubmissionsResponse struct {
	Kind string `json:"kind"`
	Data struct {
		After     string          `json:"after"`
		Dist      int             `json:"dist"`
		Modhash   string          `json:"modhash"`
		GeoFilter any             `json:"geo_filter"`
		Children  []SubmissionRes `json:"children"`
		Before    any             `json:"before"`
	} `json:"data"`
}
