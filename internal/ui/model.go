package ui

import (
	"context"
	"errors"
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
	screenSelect  = "select"
	screenSudo    = "sudo"
	screenRun     = "running"
	screenDone    = "done"
	screenSummary = "summary"
)

type model struct {
	screen   string
	modeIdx  int
	tasks    []tasks.Task
	current  int
	checking bool

	startTime time.Time
	endTime   time.Time

	sudoMode string
	sudoErr  string
	pwInput  textinput.Model

	spinner  spinner.Model
	progress progress.Model

	// failedTasks / skippedTasks track indices that did not complete cleanly.
	// A run finishes on screenSummary when either is non-empty, otherwise
	// screenDone.
	failedTasks  []int
	skippedTasks []int
	// lastErr holds the most recent task error text for the summary view.
	lastErr string
}

type taskDoneMsg struct{ idx int }
type taskErrorMsg struct {
	idx int
	err error
	out string
}
type taskSkippedMsg struct {
	idx int
	msg string
}
type sudoCheckMsg struct{ needed bool }
type sudoResultMsg struct {
	ok  bool
	err error
}
type sudoRefreshMsg struct{}

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
		return m.advance(msg.idx)

	case taskSkippedMsg:
		m.tasks[msg.idx].Status = tasks.StatusSkipped
		m.skippedTasks = append(m.skippedTasks, msg.idx)
		return m.advance(msg.idx)

	case taskErrorMsg:
		m.tasks[msg.idx].Status = tasks.StatusError
		m.failedTasks = append(m.failedTasks, msg.idx)
		m.lastErr = formatTaskError(msg.err, msg.out)
		return m.advance(msg.idx)

	case sudoRefreshMsg:
		// Keep the sudo timestamp warm during long runs so a later sudo task
		// does not hit a stale credential. Purely a side effect; the result
		// is intentionally ignored — a real failure surfaces as a task error.
		if m.screen == screenRun {
			return m, scheduleSudoRefresh()
		}
		return m, nil
	}
	return m, nil
}

// advance moves to the next runnable task or finishes the run. Tasks already
// in a terminal state (Done/Skipped/Error) are skipped, which lets a retry
// re-run only the previously failed tasks without touching the rest.
func (m model) advance(idx int) (model, tea.Cmd) {
	next := idx + 1
	for next < len(m.tasks) && m.tasks[next].Status != tasks.StatusPending {
		next++
	}
	if next >= len(m.tasks) {
		m.endTime = time.Now()
		if len(m.failedTasks) == 0 && len(m.skippedTasks) == 0 {
			m.screen = screenDone
			return m, tea.Quit
		}
		m.screen = screenSummary
		return m, m.progress.SetPercent(1.0)
	}
	pct := float64(next) / float64(len(m.tasks))
	m.current = next
	m.tasks[next].Status = tasks.StatusRunning
	return m, tea.Batch(m.runTask(next), m.progress.SetPercent(pct))
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
		if m.screen == screenDone || m.screen == screenSummary {
			return m, tea.Quit
		}
	case "r":
		if m.screen == screenSummary && len(m.failedTasks) > 0 {
			return m.retryFailed()
		}
	default:
		if m.screen == screenDone || m.screen == screenSummary {
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
	m.failedTasks = nil
	m.skippedTasks = nil
	m.lastErr = ""
	m.tasks[0].Status = tasks.StatusRunning
	return tea.Batch(
		m.runTask(0),
		func() tea.Msg { return m.spinner.Tick() },
		scheduleSudoRefresh(),
	)
}

// retryFailed resets every failed task to Pending and re-runs from the first
// one. advance() skips already-resolved tasks, so only the failed tasks
// actually execute again.
func (m model) retryFailed() (model, tea.Cmd) {
	retry := append([]int(nil), m.failedTasks...)
	m.failedTasks = nil
	m.lastErr = ""
	for _, i := range retry {
		m.tasks[i].Status = tasks.StatusPending
	}
	if len(retry) == 0 {
		return m, nil
	}
	first := retry[0]
	m.current = first
	m.tasks[first].Status = tasks.StatusRunning
	m.screen = screenRun
	m.startTime = time.Now()
	return m, tea.Batch(m.runTask(first), func() tea.Msg { return m.spinner.Tick() })
}

func checkSudo() tea.Cmd {
	return func() tea.Msg {
		err := exec.Command("sudo", "-n", "true").Run()
		return sudoCheckMsg{needed: err != nil}
	}
}

// scheduleSudoRefresh fires a silent sudo -n true after the refresh interval
// to keep the sudo credential timestamp warm. See sudoRefreshMsg.
func scheduleSudoRefresh() tea.Cmd {
	return tea.Tick(sudoRefreshInterval, func(time.Time) tea.Msg {
		_ = exec.Command("sudo", "-n", "true").Run()
		return sudoRefreshMsg{}
	})
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

// taskTimeout picks a generous deadline per command type. apt work can take
// several minutes during a big upgrade; everything else is capped at a minute
// so a hung command never freezes the TUI indefinitely.
func taskTimeout(command string) time.Duration {
	switch {
	case strings.Contains(command, "apt "),
		strings.Contains(command, "npm "),
		strings.Contains(command, "pnpm "):
		return 5 * time.Minute
	}
	return 60 * time.Second
}

// sudoRefreshInterval keeps the sudo timestamp fresh between tasks without
// spamming. Must be well under the typical sudo timestamp_timeout (5 min).
const sudoRefreshInterval = 60 * time.Second

func (m model) runTask(idx int) tea.Cmd {
	t := m.tasks[idx]
	return func() tea.Msg {
		bin := t.Require
		if bin == "" {
			bin = tasks.ExtractBinary(t.Command)
		}
		if bin != "" {
			if _, err := exec.LookPath(bin); err != nil {
				return taskSkippedMsg{idx: idx, msg: bin + " not installed"}
			}
		}
		timeout := taskTimeout(t.Command)
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		cmd := exec.CommandContext(ctx, "sh", "-c", t.Command)
		out, err := cmd.CombinedOutput()
		if ctx.Err() == context.DeadlineExceeded {
			return taskErrorMsg{idx: idx, err: errors.New("timed out after " + timeout.String()), out: string(out)}
		}
		if err != nil {
			return taskErrorMsg{idx: idx, err: err, out: string(out)}
		}
		return taskDoneMsg{idx: idx}
	}
}

// formatTaskError strips apt's boilerplate "stable CLI interface" warning and
// annotates dpkg-lock contention with an actionable hint. Returns a single
// trimmed string suitable for the summary view.
func formatTaskError(err error, out string) string {
	var kept []string
	for _, l := range strings.Split(strings.TrimSpace(out), "\n") {
		if l == "" {
			continue
		}
		if strings.Contains(l, "stable CLI interface") || strings.Contains(l, "Use with caution in scripts") {
			continue
		}
		kept = append(kept, l)
	}
	body := strings.Join(kept, "\n")

	if strings.Contains(out, "lock-frontend") || strings.Contains(out, "Could not get lock") {
		hint := "Another package manager may be running (often unattended-upgrades).\n" +
			"Stop it: sudo systemctl stop unattended-upgrades"
		if body == "" {
			return err.Error() + "\n" + hint
		}
		return err.Error() + "\n" + body + "\n" + hint
	}

	if body == "" {
		return err.Error()
	}
	return err.Error() + "\n" + body
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
	case screenSummary:
		s = m.summaryView()
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
		{"🌿", "Shallow Clean", "Fast daily cleanup (~20s)",
			fmt.Sprintf("%d tasks · safe for everyday use", len(tasks.ShallowTasks()))},
		{"🔥", "Deep Clean", "Full system cleanup (~90s)",
			fmt.Sprintf("%d tasks · requires sudo", len(tasks.DeepTasks()))},
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

// renderTaskRow renders a single task line for run/summary views. The
// spinner view is passed in so callers share one spinner instance.
func renderTaskRow(t tasks.Task, spinnerView string) string {
	var mark, name string
	switch t.Status {
	case tasks.StatusDone:
		mark = doneStyle.Render("✅")
		name = doneStyle.Render(t.Name)
	case tasks.StatusRunning:
		mark = runStyle.Render(spinnerView)
		name = runStyle.Render(t.Name)
	case tasks.StatusError:
		mark = errStyle.Render("❌")
		name = errStyle.Render(t.Name)
	case tasks.StatusSkipped:
		mark = pendStyle.Render("⊘")
		name = pendStyle.Render(t.Name + "  (skipped)")
	default:
		mark = pendStyle.Render("·")
		name = pendStyle.Render(t.Name)
	}
	return fmt.Sprintf("  %s  %s\n", mark, name)
}

func (m model) runView() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("🧹  Linux Cleanup CLI") + "  " +
		modeBadgeStyle.Render("["+tasks.ModeName(m.modeIdx)+"]"))
	b.WriteString("\n\n")

	for _, t := range m.tasks {
		b.WriteString(renderTaskRow(t, m.spinner.View()))
	}

	n := len(m.tasks)
	done := countDone(m.tasks)
	b.WriteString("\n")
	b.WriteString(m.progress.View() + "  " + pendStyle.Render(fmt.Sprintf("%d/%d tasks", done, n)))
	return b.String()
}

// successBody returns the one-line summary shown on the done screen and in
// the goodbye card printed after exit.
func (m model) successBody() string {
	count := len(m.tasks)
	elapsed := m.endTime.Sub(m.startTime).Round(time.Second)
	return fmt.Sprintf("%d tasks completed in %s", count, elapsed)
}

func (m model) doneView() string {
	return SuccessCard("Done", m.successBody()) + "\n"
}

func (m model) summaryView() string {
	var b strings.Builder
	for _, t := range m.tasks {
		b.WriteString(renderTaskRow(t, m.spinner.View()))
	}
	b.WriteString("\n")

	done := countDone(m.tasks)
	failed := len(m.failedTasks)
	skipped := len(m.skippedTasks)
	elapsed := m.endTime.Sub(m.startTime).Round(time.Second)
	b.WriteString(modeBadgeStyle.Render(fmt.Sprintf(
		"%d done · %d failed · %d skipped · %s", done, failed, skipped, elapsed)))
	b.WriteString("\n\n")

	if m.lastErr != "" {
		b.WriteString(errStyle.Render(m.lastErr))
		b.WriteString("\n\n")
	}

	if failed > 0 {
		b.WriteString(hintStyle.Render("r retry failed   q quit"))
	} else {
		b.WriteString(hintStyle.Render("press any key to exit"))
	}
	return WarningCard("Cleanup Summary", b.String()) + "\n"
}

// finalMessage returns the card to print after the program exits so the
// result lingers in the terminal scrollback. Returns "" for states that
// should leave the terminal untouched (early quit at select/sudo screens).
func (m model) finalMessage() string {
	switch m.screen {
	case screenDone:
		return SuccessCard("Done", m.successBody())
	case screenSummary:
		return m.summaryFinalCard()
	}
	return ""
}

// summaryFinalCard renders the summary without the interactive hints, since
// those are meaningless once the program has exited.
func (m model) summaryFinalCard() string {
	var b strings.Builder
	for _, t := range m.tasks {
		b.WriteString(renderTaskRow(t, ""))
	}
	b.WriteString("\n")
	done := countDone(m.tasks)
	failed := len(m.failedTasks)
	skipped := len(m.skippedTasks)
	elapsed := m.endTime.Sub(m.startTime).Round(time.Second)
	b.WriteString(modeBadgeStyle.Render(fmt.Sprintf(
		"%d done · %d failed · %d skipped · %s", done, failed, skipped, elapsed)))
	if m.lastErr != "" {
		b.WriteString("\n\n")
		b.WriteString(errStyle.Render(m.lastErr))
	}
	return WarningCard("Cleanup Summary", b.String())
}

// FinalMessage is the exported goodbye-card accessor used by main after the
// program returns. It type-asserts the final tea.Model back to the concrete
// model and returns its finalMessage, or "" if the assertion fails.
func FinalMessage(m tea.Model) string {
	if mm, ok := m.(model); ok {
		return mm.finalMessage()
	}
	return ""
}
