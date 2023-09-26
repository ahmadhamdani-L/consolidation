package model

import "time"

type JsonData struct {
	ID        int
	FileLoc   string
	CompanyID int
	UserID    int
	Timestamp *time.Time
	Data      string
	Name      string
	Filter    struct {
		Period   string
		Versions int
	}
}

type Message struct {
	NotificationID int        `json:"id"`
	Name           string     `json:"name"`
	Company        string     `json:"company"`
	User           string     `json:"user"`
	Period         string     `json:"period"`
	Versions       int        `json:"versions"`
	Message        string     `json:"message"`
	FileUrl        string     `json:"file_url,omitempty"`
	Timestamp      *time.Time `json:"timestamp"`
}
