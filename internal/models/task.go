package models

import "time"

type Task struct {
	ID        int        `json:"id"`
	Title     string     `json:"title"`
	Completed bool       `json:"completed"`
	Deadline  *time.Time `json:"deadline,omitempty"`
	IsOverdue bool       `json:"is_overdue"`
}

type TaskRequest struct {
	Title    string     `json:"title"`
	Deadline *time.Time `json:"deadline,omitempty"`
}
