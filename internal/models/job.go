package models

import "time"

type Job struct {
	JobID          string    `db:"job_id" json:"job_id"`
	FromDate       time.Time `db:"from_date" json:"from"`
	ToDate         time.Time `db:"to_date" json:"to"`
	Status         string    `db:"status" json:"status"`
	Progress       float64   `db:"progress" json:"progress"`
	ProcessedCount int       `db:"processed_count" json:"processed"`
	TotalCount     int       `db:"total_count" json:"total"`
	ResultPath     *string   `db:"result_path" json:"download_url,omitempty"`
}

type RequestJob struct {
	From string `json:"from" binding:"required"`
	To   string `json:"to" binding:"required"`
}
