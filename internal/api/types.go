// Defines data structures that match the Clockify API responses
package api

// UserInfo contains information about the current user
type UserInfo struct {
	ID               string `json:"id"`
	DefaultWorkspace string `json:"defaultWorkspace"`
}

// Project represents a Clockify project
// The json tags tell Go how to map JSON fields to struct fields
type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	ClientID string `json:"clientId"`
	ClientName string `json:"clientName"`
}

// TimeEntryRequest is the payload we send when creating a time entry
type TimeEntryRequest struct {
	Start       string `json:"start"`       // RFC3339 format timestamp
	End         string `json:"end"`         // RFC3339 format timestamp
	ProjectID   string `json:"projectId"`   // ID of the project
	Description string `json:"description"` // Task description
}

// TimeEntryResponse represents a time entry returned from the API
// We use this to parse previous entries and extract task descriptions
type TimeEntryResponse struct {
	Description string `json:"description"`
}
