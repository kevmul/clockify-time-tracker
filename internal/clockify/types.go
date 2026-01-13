package clockify

// ViewState represents different views in the application
type ViewState int

const (
	ViewDashboard ViewState = iota
	ViewTimeEntry
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

type CustomFieldValue struct {
	CustomFieldId string `json:"customFieldId"`
	Name          string `json:"name"`
	TimeEntryId   string `json:"timeEntryId"`
	Type          string `json:"fieldType"`
	Value         string `json:"value"`
}

type Entry struct {
	Billable bool `json:"billable"`
	CostRate struct {
		Amount   int32  `json:"amount"`
		Currency string `json:"currency"`
	}
	CustomFieldValues []CustomFieldValue
	Description       string `json:"description"`
	HourlyRate        struct {
		Amount   int32  `json:"amount"`
		Currency string `json:"currency"`
	}
	Id           string   `json:"id"`
	IsLocked     bool     `json:"isLocked"`
	KioskId      string   `json:"kioskId"`
	ProjectId    string   `json:"projectId"`
	TagIds       []string `json:"tagIds"`
	TaskId       string   `json:"taskId"`
	TimeInterval struct {
		Duration string `json:"duration"`
		End      string `json:"end"`
		Start    string `json:"start"`
	}
	Type        string `json:"type"`
	UserId      string `json:"userId"`
	WorkspaceId string `json:"workspaceId"`
}
