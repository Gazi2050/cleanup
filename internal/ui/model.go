package ui

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"charm.land/bubbles/v2/progress"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"

	"github.com/Gazi2050/cleanup/internal/tasks"
)

const (
	screenSelect = "select"
	screenSudo   = "sudo"
	screenRun    = "running"
	screenDone   = "done"
	screenError  = "error"
)

type model struct {
	screen   string
	modeIdx  int
	tasks    []tasks.Task
	current  int
	checking bool

	startTime time.Time
	endTime   time.Time
	errMsg    string

	sudoMode string
	sudoErr  string
	pwInput  textinput.Model

	spinner  spinner.Model
	progress progress.Model
}

type taskDoneMsg struct{ idx int }
type taskErrorMsg struct {
	idx int
	err error
	out string
}
type sudoCheckMsg struct{ needed bool }
type sudoResultMsg struct {
	ok  bool
	err error
}
type toastTimeoutMsg struct{}

func InitialModel() model {
	s := spinner.New(
		spinner.WithSpinner(spinner.MiniDot),
		spinner.WithStyle(runStyle),
	)
	p := progress.New(
		progress.WithoutPercentage(),
		progress.WithScaled(true),
		progress.WithColors(theme.Header, theme.Mode),
		progress.WithWidth(30),
	)
	ti := textinput.New()
	ti.Prompt = "Password: "
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '*'
	ti.CharLimit = 128
	return model{
		screen:   screenSelect,
		spinner:  s,
		progress: p,
		pwInput:  ti,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		SetTerminalWidth(msg.Width)
		return m, nil

	case tea.KeyPressMsg:
		if m.screen == screenSudo {
			return m.updateSudoKeys(msg)
		}
		return m.updateGlobalKeys(msg)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case progress.FrameMsg:
		var cmd tea.Cmd
		m.progress, cmd = m.progress.Update(msg)
		return m, cmd

	case sudoCheckMsg:
		m.checking = false
		if msg.needed {
			m.screen = screenSudo
			return m, m.pwInput.Focus()
		}
		return m, m.startRunning()

	case sudoResultMsg:
		if msg.ok {
			m.sudoErr = ""
			m.pwInput.SetValue("")
			return m, m.startRunning()
		}
		m.sudoErr = "authentication failed"
		m.pwInput.SetValue("")
		return m, m.pwInput.Focus()

	case taskDoneMsg:
		m.tasks[msg.idx].Status = tasks.StatusDone
		next := msg.idx + 1
		pct := float64(next) / float64(len(m.tasks))
		if next >= len(m.tasks) {
			m.endTime = time.Now()
			m.screen = screenDone
			return m, tea.Batch(
				m.progress.SetPercent(1.0),
				tea.Tick(3*time.Second, func(time.Time) tea.Msg { return toastTimeoutMsg{} }),
			)
		}
		m.current = next
		m.tasks[next].Status = tasks.StatusRunning
		return m, tea.Batch(m.runTask(next), m.progress.SetPercent(pct))

	case taskErrorMsg:
		m.tasks[msg.idx].Status = tasks.StatusError
		m.endTime = time.Now()
		m.screen = screenError
		if msg.out != "" {
			m.errMsg = strings.TrimSpace(msg.err.Error() + "\n" + msg.out)
		} else {
			m.errMsg = msg.err.Error()
		}
		return m, nil

	case toastTimeoutMsg:
		return m, tea.Quit
	}
	return m, nil
}

func (m model) updateGlobalKeys(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "up", "k":
		if m.screen == screenSelect && !m.checking && m.modeIdx > 0 {
			m.modeIdx--
		}
	case "down", "j":
		if m.screen == screenSelect && !m.checking && m.modeIdx < 1 {
			m.modeIdx++
		}
	case "enter":
		if m.screen == screenSelect && !m.checking {
			if m.modeIdx == 0 {
				m.tasks = tasks.ShallowTasks()
				m.sudoMode = "shallow"
			} else {
				m.tasks = tasks.DeepTasks()
				m.sudoMode = "deep"
			}
			m.checking = true
			return m, tea.Batch(checkSudo(), func() tea.Msg { return m.spinner.Tick() })
		}
	default:
		if m.screen == screenDone || m.screen == screenError {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) updateSudoKeys(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "enter":
		pw := m.pwInput.Value()
		m.pwInput.SetValue("")
		m.sudoErr = ""
		return m, validateSudo(pw)
	default:
		var cmd tea.Cmd
		m.pwInput, cmd = m.pwInput.Update(msg)
		return m, cmd
	}
}

func (m *model) startRunning() tea.Cmd {
	m.screen = screenRun
	m.startTime = time.Now()
	m.current = 0
	m.tasks[0].Status = tasks.StatusRunning
	return tea.Batch(m.runTask(0), func() tea.Msg { return m.spinner.Tick() })
}

func checkSudo() tea.Cmd {
	return func() tea.Msg {
		err := exec.Command("sudo", "-n", "true").Run()
		return sudoCheckMsg{needed: err != nil}
	}
}

func countDone(list []tasks.Task) int {
	n := 0
	for _, t := range list {
		if t.Status == tasks.StatusDone {
			n++
		}
	}
	return n
}

func validateSudo(pw string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("sudo", "-S", "-v")
		cmd.Stdin = strings.NewReader(pw + "\n")
		err := cmd.Run()
		return sudoResultMsg{ok: err == nil, err: err}
	}
}

func (m model) runTask(idx int) tea.Cmd {
	task := m.tasks[idx]
	return func() tea.Msg {
		cmd := exec.Command("sh", "-c", task.Command)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return taskErrorMsg{idx: idx, err: err, out: string(out)}
		}
		return taskDoneMsg{idx: idx}
	}
}

func (m model) View() tea.View {
	var s string
	switch m.screen {
	case screenSelect:
		s = m.selectView()
	case screenSudo:
		s = m.sudoView()
	case screenRun:
		s = m.runView()
	case screenDone:
		s = m.doneView()
	case screenError:
		s = m.errorView()
	}
	v := tea.NewView(s)
	v.AltScreen = true
	return v
}

func (m model) selectView() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("🧹  Linux Cleanup CLI"))
	b.WriteString("\n\n")

	if m.checking {
		b.WriteString(runStyle.Render(m.spinner.View()) + "  Checking sudo access...")
		b.WriteString("\n\n")
		b.WriteString(hintStyle.Render("please wait"))
		return b.String()
	}

	b.WriteString("Select cleanup mode:\n\n")
	modes := []struct {
		icon string
		name string
		desc string
		meta string
	}{
		{"🌿", "Shallow Clean", "Fast daily cleanup (~20s)", "5 tasks · safe for everyday use"},
		{"🔥", "Deep Clean", "Full system cleanup (~90s)", "11 tasks · requires sudo"},
	}
	for i, mode := range modes {
		marker := pendStyle.Render("( )")
		if i == m.modeIdx {
			marker = cursorStyle.Render("(•)")
		}
		b.WriteString(fmt.Sprintf("  %s %s\n", marker, cardTitleStyle.Render(mode.icon+"  "+mode.name)))
		b.WriteString(fmt.Sprintf("      %s\n", cardDescStyle.Render(mode.desc)))
		b.WriteString(fmt.Sprintf("      %s\n", cardMetaStyle.Render(mode.meta)))
		b.WriteString("\n")
	}
	b.WriteString(hintStyle.Render("↑↓ navigate   enter select   q quit"))
	return b.String()
}

func (m model) sudoView() string {
	var b strings.Builder
	b.WriteString(m.pwInput.View())
	if m.sudoErr != "" {
		b.WriteString("\n")
		b.WriteString(pwErrStyle.Render(m.sudoErr))
	}
	b.WriteString("\n\n")
	b.WriteString(hintStyle.Render("enter submit · ctrl+c cancel"))
	return RenderBox(BoxWarning, "🔒  Sudo password required", b.String())
}

func (m model) runView() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("🧹  Linux Cleanup CLI") + "  " +
		modeBadgeStyle.Render("["+tasks.ModeName(m.modeIdx)+"]"))
	b.WriteString("\n\n")

	for _, t := range m.tasks {
		var mark, name string
		switch t.Status {
		case tasks.StatusDone:
			mark = doneStyle.Render("✅")
			name = doneStyle.Render(t.Name)
		case tasks.StatusRunning:
			mark = runStyle.Render(m.spinner.View())
			name = runStyle.Render(t.Name)
		case tasks.StatusError:
			mark = errStyle.Render("❌")
			name = errStyle.Render(t.Name)
		default:
			mark = pendStyle.Render("·")
			name = pendStyle.Render(t.Name)
		}
		b.WriteString(fmt.Sprintf("  %s  %s\n", mark, name))
	}

	n := len(m.tasks)
	done := countDone(m.tasks)
	b.WriteString("\n")
	b.WriteString(m.progress.View() + "  " + pendStyle.Render(fmt.Sprintf("%d/%d tasks", done, n)))
	return b.String()
}

func (m model) doneView() string {
	count := len(m.tasks)
	elapsed := m.endTime.Sub(m.startTime).Round(time.Second)
	body := fmt.Sprintf("%d tasks completed in %s", count, elapsed)
	return SuccessCard("Done", body) + "\n"
}

func (m model) errorView() string {
	var b strings.Builder
	for _, t := range m.tasks {
		var mark, name string
		switch t.Status {
		case tasks.StatusDone:
			mark = doneStyle.Render("✅")
			name = doneStyle.Render(t.Name)
		case tasks.StatusError:
			mark = errStyle.Render("❌")
			name = errStyle.Render(t.Name + "  ← failed here")
		default:
			mark = pendStyle.Render("·")
			name = pendStyle.Render(t.Name)
		}
		b.WriteString(fmt.Sprintf("  %s  %s\n", mark, name))
	}
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("Error: %s\n", errStyle.Render(m.errMsg)))
	b.WriteString(hintStyle.Render("Tip: ensure sudo is available and disk is writable"))
	card := ErrorCard("Task Failed", b.String())
	return card + "\n" + hintStyle.Render("press any key to exit")
}
