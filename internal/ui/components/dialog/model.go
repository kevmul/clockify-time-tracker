package dialog

type Modal struct {
	Title   string
	Content string
	Width   int
	Height  int
	Visible bool
}

func NewModal(title string) Modal {
	return Modal{
		Title:   title,
		Width:   60,
		Height:  20,
		Visible: false,
	}
}

func (m *Modal) Show(content string) {
	m.Content = content
	m.Visible = true
}

func (m *Modal) Hide() {
	m.Visible = false
}
