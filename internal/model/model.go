package model

import "time"

type APOD struct {
	ID             int       `json:"id"`
	Copyright      string    `json:"copyright"`
	Date           string    `json:"date"`
	ParsedDate     time.Time `json:"-"`
	Explanation    string    `json:"explanation"`
	MediaType      string    `json:"media_type"`
	ServiceVersion string    `json:"service_version"`
	Title          string    `json:"title"`
	URL            string    `json:"url"`
}
