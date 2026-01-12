// Defines the different steps/screens in our application flow
package timeform

// These constants represent each screen in our UI flow
// Using constants (instead of magic numbers) makes the code more readable
const (
	stepDateSelect    = iota // 0 - Select which date to log time for
	stepProjectSelect        // 1 - Select which project
	stepTimeInput            // 2 - Enter time range (e.g., "9a - 5p")
	stepTaskInput            // 3 - Enter task description
	stepConfirm              // 4 - Review and confirm the entry
	stepComplete             // 5 - Show success message
)
