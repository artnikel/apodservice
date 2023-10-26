// Package constants - all project constants
package constants

import "time"

const (
	// DateLayout is format yyyy-mm-dd
	DateLayout = "2006-01-02"
	// MediaPath is path to save images
	MediaPath = "storage"
	// WorkFrequency is worker request frequency
	WorkFrequency = 24 * time.Hour
	// MethodGet - http method GET of request
	MethodGet = "GET"
	// ServerTimeout is read and write timeout of server config
	ServerTimeout = 10 * time.Second
	// InvalidMethod is error if method invalid
	InvalidMethod = "Invalid HTTP method"
)
