package messages

import "clockify-time-tracker/internal/clockify"

// Re-export types for backward compatibility
type NavigationMsg = clockify.NavigationMsg
type ErrMsg = clockify.ErrMsg
type ProjectsMsg = clockify.ProjectsMsg
type TasksMsg = clockify.TasksMsg
type UserInfoMsg = clockify.UserInfoMsg
type SetLoadingMsg = clockify.SetLoadingMsg
type SubmitSuccessMsg = clockify.SubmitSuccessMsg
type QuittingAppMsg = clockify.QuittingAppMsg
