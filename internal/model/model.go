// Package model contains models which represents data of some object
package model

import "time"

// APOD is a struct of metadata NASA APOD
type APOD struct {
	ID             int       `json:"-"`
	Copyright      string    `json:"copyright"`
	Date           string    `json:"date"`
	ParsedDate     time.Time `json:"-"`
	Explanation    string    `json:"explanation"`
	MediaType      string    `json:"media_type"`
	ServiceVersion string    `json:"service_version"`
	Title          string    `json:"title"`
	URL            string    `json:"url"`
}
