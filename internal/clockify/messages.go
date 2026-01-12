package clockify

type NavigationMsg struct {
	View ViewState
}

type ErrMsg error

type ProjectsMsg []Project
type TasksMsg []string
type UserInfoMsg struct {
	WorkspaceID string
	UserID      string
}
type SubmitSuccessMsg struct{}

type SetLoadingMsg struct{}

type QuittingAppMsg struct{}
