package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"charm.land/bubbles/v2/progress"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
)

const (
	screenSelect = "select"
	screenRun    = "running"
	screenDone   = "done"
	screenError  = "error"
)

type model struct {
	screen    string
	modeIdx   int
	tasks     []Task
	current   int
	startTime time.Time
	endTime   time.Time
	errMsg    string
	spinner   spinner.Model
	progress  progress.Model
}

type taskDoneMsg struct{ idx int }
type taskErrorMsg struct {
	idx int
	err error
	out string
}

func initialModel() model {
	s := spinner.New(
		spinner.WithSpinner(spinner.MiniDot),
		spinner.WithStyle(runStyle),
	)
	p := progress.New(
		progress.WithoutPercentage(),
		progress.WithScaled(true),
		progress.WithColors(colorHeader, colorMode),
	)
	return model{
		screen:   screenSelect,
		spinner:  s,
		progress: p,
	}
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg { return m.spinner.Tick() }
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.progress.SetWidth(msg.Width - 8)
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.screen == screenSelect && m.modeIdx > 0 {
				m.modeIdx--
			}
		case "down", "j":
			if m.screen == screenSelect && m.modeIdx < 1 {
				m.modeIdx++
			}
		case "enter":
			if m.screen == screenSelect {
				if m.modeIdx == 0 {
					m.tasks = ShallowTasks()
				} else {
					m.tasks = DeepTasks()
				}
				m.screen = screenRun
				m.startTime = time.Now()
				m.current = 0
				m.tasks[0].Status = StatusRunning
				return m, tea.Batch(m.runTask(0), func() tea.Msg { return m.spinner.Tick() })
			}
		default:
			if m.screen == screenDone || m.screen == screenError {
				return m, tea.Quit
			}
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case progress.FrameMsg:
		var cmd tea.Cmd
		m.progress, cmd = m.progress.Update(msg)
		return m, cmd

	case taskDoneMsg:
		m.tasks[msg.idx].Status = StatusDone
		next := msg.idx + 1
		pct := float64(next) / float64(len(m.tasks))
		if next >= len(m.tasks) {
			m.endTime = time.Now()
			m.screen = screenDone
			return m, m.progress.SetPercent(1.0)
		}
		m.current = next
		m.tasks[next].Status = StatusRunning
		return m, tea.Batch(m.runTask(next), m.progress.SetPercent(pct))

	case taskErrorMsg:
		m.tasks[msg.idx].Status = StatusError
		m.endTime = time.Now()
		m.screen = screenError
		if msg.out != "" {
			m.errMsg = strings.TrimSpace(msg.err.Error() + "\n" + msg.out)
		} else {
			m.errMsg = msg.err.Error()
		}
		return m, nil
	}
	return m, nil
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
	b.WriteString(titleStyle.Render("🧹  Ubuntu Cleanup CLI"))
	b.WriteString("\n\n")
	b.WriteString("Select cleanup mode:\n\n")

	modes := []struct {
		icon string
		name string
		desc string
	}{
		{"🌿", "Shallow Clean", "Fast daily cleanup (~20s)"},
		{"🔥", "Deep Clean", "Full system cleanup (~90s)"},
	}
	for i, mode := range modes {
		cursor := "   "
		if i == m.modeIdx {
			cursor = cursorStyle.Render("▶  ")
		}
		b.WriteString(fmt.Sprintf("%s%s  %s\n", cursor, mode.icon, mode.name))
		b.WriteString(fmt.Sprintf("       %s\n\n", mode.desc))
	}
	b.WriteString(hintStyle.Render("↑↓ navigate   enter select   q quit"))
	return boxStyle.Render(b.String())
}

func (m model) runView() string {
	var b strings.Builder
	header := titleStyle.Render("🧹  Ubuntu Cleanup CLI") + "  " + modeBadgeStyle.Render("["+modeName(m.modeIdx)+"]")
	b.WriteString(header)
	b.WriteString("\n\n")

	for _, t := range m.tasks {
		var mark, name string
		switch t.Status {
		case StatusDone:
			mark = doneStyle.Render("✅")
			name = doneStyle.Render(t.Name)
		case StatusRunning:
			mark = runStyle.Render(m.spinner.View())
			name = runStyle.Render(t.Name)
		case StatusError:
			mark = errStyle.Render("❌")
			name = errStyle.Render(t.Name)
		default:
			mark = pendStyle.Render("·")
			name = pendStyle.Render(t.Name)
		}
		b.WriteString(fmt.Sprintf("  %s  %s\n", mark, name))
	}

	count := 0
	for _, t := range m.tasks {
		if t.Status == StatusDone {
			count++
		}
	}
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  %s  %d/%d tasks\n", m.progress.View(), count, len(m.tasks)))
	return boxStyle.Render(b.String())
}

func (m model) doneView() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("✨  Cleanup Complete!"))
	b.WriteString("\n\n")
	for _, t := range m.tasks {
		b.WriteString(fmt.Sprintf("  %s  %s\n", doneStyle.Render("✅"), doneStyle.Render(t.Name)))
	}
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  %s  %d/%d ✓\n", m.progress.View(), len(m.tasks), len(m.tasks)))
	b.WriteString("\n")
	elapsed := m.endTime.Sub(m.startTime).Round(time.Second)
	b.WriteString(fmt.Sprintf("🎉  Your system is fresh and clean!\n     Time taken: %s\n\n", elapsed))
	b.WriteString(hintStyle.Render("press any key to exit"))
	return boxStyle.Render(b.String())
}

func (m model) errorView() string {
	var b strings.Builder
	b.WriteString(errStyle.Render("⚠️   Task Failed"))
	b.WriteString("\n\n")
	for _, t := range m.tasks {
		var mark, name string
		switch t.Status {
		case StatusDone:
			mark = doneStyle.Render("✅")
			name = doneStyle.Render(t.Name)
		case StatusError:
			mark = errStyle.Render("❌")
			name = errStyle.Render(t.Name + "  ← failed here")
		default:
			mark = pendStyle.Render("·")
			name = pendStyle.Render(t.Name)
		}
		b.WriteString(fmt.Sprintf("  %s  %s\n", mark, name))
	}
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  Error: %s\n", errStyle.Render(m.errMsg)))
	b.WriteString(fmt.Sprintf("  %s\n", hintStyle.Render("Tip: ensure passwordless sudo is configured (visudo)")))
	b.WriteString("\n")
	b.WriteString(hintStyle.Render("press any key to exit"))
	return boxStyle.Render(b.String())
}
