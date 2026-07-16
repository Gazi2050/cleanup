package tasks

import (
	"fmt"
	"strings"
)

type TaskStatus int

const (
	StatusPending TaskStatus = iota
	StatusRunning
	StatusDone
	StatusError
	StatusSkipped
)

// Task represents a single cleanup step.
//
// Require names a binary that must exist on PATH before the task runs. If it
// is missing the task is marked StatusSkipped instead of being executed.
// When empty, runTask auto-detects the binary from Command via
// tasks.ExtractBinary (sudo/nohup stripped, first token used).
type Task struct {
	Name    string
	Command string
	Status  TaskStatus
	Require string
}

// aptBase holds the non-interactive flags every apt invocation needs.
//
// DPkg::Lock::Timeout makes apt wait up to 5 minutes for the dpkg frontend
// lock instead of failing instantly when unattended-upgrades (or another
// package manager) holds it. force-confold/force-confdef keep config prompts
// from blocking an unattended upgrade.
const aptBase = "sudo apt -y " +
	"-o DPkg::Lock::Timeout=300 " +
	"-o Dpkg::Options::=--force-confold " +
	"-o Dpkg::Options::=--force-confdef "

// aptCmd wraps an apt subcommand with the non-interactive base flags.
func aptCmd(args string) string { return aptBase + args }

// ExtractBinary returns the primary executable name from a shell command
// string, or "" if none can be determined. Strips leading "sudo"/"nohup"
// wrappers and "env VAR=val" assignments so e.g. "sudo apt update" yields
// "apt". Compound commands (&& or ;) are inspected only at the first
// segment, since the first binary is the task's primary tool.
func ExtractBinary(command string) string {
	fields := strings.Fields(command)
	for len(fields) > 0 && (fields[0] == "sudo" || fields[0] == "nohup") {
		fields = fields[1:]
	}
	if len(fields) == 0 {
		return ""
	}
	if fields[0] == "env" {
		i := 1
		for i < len(fields) && strings.Contains(fields[i], "=") {
			i++
		}
		if i < len(fields) {
			return fields[i]
		}
		return ""
	}
	return fields[0]
}

// tmpKeepNames lists top-level /tmp entries that must never be deleted because
// they are live IPC sockets, locks, or per-service dirs held by the running
// session. Deleting them mid-session breaks X11, SSH agents, PulseAudio,
// systemd private dirs, and snap daemons.
var tmpKeepNames = []string{
	"X11-unix",
	".X*-lock",
	".ICE-unix",
	".font-unix",
	"ssh-*",
	"pulse-*",
	"systemd-*",
	"snap-*",
}

// tmpClearCmd clears /tmp and /var/tmp while preserving the live entries in
// tmpKeepNames. The -prune/-o -exec idiom skips kept names and removes
// everything else at the top level only (no recursion into subdirs of kept
// paths). Errors from busy files are swallowed by "2>/dev/null || true".
func tmpClearCmd() string {
	parts := make([]string, 0, len(tmpKeepNames))
	for _, n := range tmpKeepNames {
		parts = append(parts, fmt.Sprintf("-name %q", n))
	}
	prune := "\\( " + strings.Join(parts, " -o ") + " \\) -prune"
	return "sudo find /tmp /var/tmp -mindepth 1 -maxdepth 1 " + prune +
		" -o -exec rm -rf {} + 2>/dev/null || true"
}

func ShallowTasks() []Task {
	return []Task{
		{Name: "Clear Trash", Command: "rm -rf ~/.local/share/Trash/*"},
		{Name: "Clear User Cache", Command: "rm -rf ~/.cache/*"},
		{Name: "Remove .tmp files", Command: "find /tmp ~/.cache -name \"*.tmp\" -delete 2>/dev/null || true"},
		{Name: "Vacuum journals (3 days)", Command: "sudo journalctl --vacuum-time=3d"},
		{Name: "APT autoclean", Command: aptCmd("autoclean")},
	}
}

func DeepTasks() []Task {
	return []Task{
		{Name: "APT update", Command: aptCmd("update")},
		{Name: "APT full-upgrade", Command: aptCmd("full-upgrade")},
		{Name: "Update global npm packages", Command: "npm update -g", Require: "npm"},
		{Name: "Update global pnpm packages", Command: "pnpm update -g --latest", Require: "pnpm"},
		{Name: "APT autoremove", Command: aptCmd("autoremove")},
		{Name: "APT autoclean + clean", Command: aptCmd("autoclean && sudo apt clean")},
		{Name: "Clean npm cache", Command: "npm cache clean --force", Require: "npm"},
		{Name: "Prune pnpm store", Command: "pnpm store prune", Require: "pnpm"},
		{Name: "Clear /tmp and /var/tmp", Command: tmpClearCmd()},
		{Name: "Clear user cache + trash", Command: "rm -rf ~/.cache/* ~/.local/share/Trash/*"},
		{Name: "Remove .tmp files", Command: "find /tmp ~/.cache ~/.local -name \"*.tmp\" -delete 2>/dev/null || true"},
		{Name: "Vacuum journals (3 days)", Command: "sudo journalctl --vacuum-time=3d"},
	}
}

func ModeName(idx int) string {
	if idx == 0 {
		return "SHALLOW CLEAN"
	}
	return "DEEP CLEAN"
}
