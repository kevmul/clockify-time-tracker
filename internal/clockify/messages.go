package clockify

type NavigationMsg struct {
	View ViewState
}

type ErrMsg error

type ProjectsMsg []Project
type UserInfoMsg struct {
	WorkspaceID string
	UserID      string
}

type SetLoadingMsg struct{}

type QuittingAppMsg struct{}

// Time Entries
type SubmitSuccessMsg struct{}
type EntriesMsg struct {
	Entries  []Entry
	Projects []Project
}

type CreateOrEdit int

const (
	Create CreateOrEdit = iota
	Edit
)

type CreateOrEditEntryMsg struct {
	Type CreateOrEdit
}

// Tasks
type TasksMsg []string
