package types

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

// Message types
type NavigationMsg struct {
	Item  string
	Index int
}

type ErrMsg error

type ProjectsMsg []Project
type TasksMsg []string
type UserInfoMsg struct {
	WorkspaceID string
	UserID      string
}
type SubmitSuccessMsg struct{}
