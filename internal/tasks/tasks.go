package tasks

type TaskStatus int

const (
	StatusPending TaskStatus = iota
	StatusRunning
	StatusDone
	StatusError
)

type Task struct {
	Name    string
	Command string
	Status  TaskStatus
}

func ShallowTasks() []Task {
	return []Task{
		{Name: "Clear Trash", Command: "rm -rf ~/.local/share/Trash/*"},
		{Name: "Clear User Cache", Command: "rm -rf ~/.cache/*"},
		{Name: "Remove .tmp files", Command: "find /tmp ~/.cache -name \"*.tmp\" -delete 2>/dev/null || true"},
		{Name: "Vacuum journals (3 days)", Command: "sudo journalctl --vacuum-time=3d"},
		{Name: "APT autoclean", Command: "sudo apt autoclean"},
	}
}

func DeepTasks() []Task {
	return []Task{
		{Name: "APT update", Command: "sudo apt update"},
		{Name: "APT upgrade", Command: "sudo apt upgrade -y"},
		{Name: "APT full-upgrade", Command: "sudo apt full-upgrade -y"},
		{Name: "APT autoremove", Command: "sudo apt autoremove -y"},
		{Name: "APT autoclean + clean", Command: "sudo apt autoclean && sudo apt clean"},
		{Name: "Clear /tmp and /var/tmp", Command: "sudo rm -rf /tmp/* /var/tmp/*"},
		{Name: "Clear user cache + trash", Command: "rm -rf ~/.cache/* ~/.local/share/Trash/*"},
		{Name: "Remove .tmp files", Command: "find /tmp ~/.cache ~/.local -name \"*.tmp\" -delete 2>/dev/null || true"},
		{Name: "Clean npm cache", Command: "npm cache clean --force && rm -rf ~/.npm"},
		{Name: "Prune pnpm store", Command: "pnpm store prune && rm -rf ~/.pnpm-store"},
		{Name: "Vacuum journals (3 days)", Command: "sudo journalctl --vacuum-time=3d"},
	}
}

func ModeName(idx int) string {
	if idx == 0 {
		return "SHALLOW CLEAN"
	}
	return "DEEP CLEAN"
}
