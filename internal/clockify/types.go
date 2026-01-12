package clockify

// ViewState represents different views in the application
type ViewState int

const (
	ViewDashboard ViewState = iota
	ViewTimeList
	ViewSettings
)

// API types
type UserInfo struct {
	ID               string `json:"id"`
	DefaultWorkspace string `json:"defaultWorkspace"`
}

type Project struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ClientID   string `json:"clientId"`
	ClientName string `json:"clientName"`
}

type TimeEntryRequest struct {
	Start       string `json:"start"`
	End         string `json:"end"`
	ProjectID   string `json:"projectId"`
	Description string `json:"description"`
}

type TimeEntryResponse struct {
	Description string `json:"description"`
}
