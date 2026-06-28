# cleanup

Interactive TUI for Ubuntu system cleanup. Built with Go + Bubble Tea v2.

## Modes

| Mode | Tasks | Time | Sudo |
|------|-------|------|------|
| Shallow | 5 | ~20s | partial (journalctl, apt) |
| Deep | 12 | ~90s | yes |

## Prerequisites

- Ubuntu (or apt-based distro)
- Go 1.26+ (only for building)
- **Passwordless sudo configured** for the current user

### Configure passwordless sudo

```bash
sudo visudo -f /etc/sudoers.d/cleanup
```

Add (replace `YOURUSER` with your username):

```
YOURUSER ALL=(ALL) NOPASSWD: /usr/bin/journalctl, /usr/bin/apt, /usr/bin/rm
```

Validate: `sudo -k && sudo -n true` should print nothing (exit 0).

## Build

```bash
go build -ldflags="-s -w" -o cleanup .
```

## Install globally

```bash
sudo mv cleanup /usr/local/bin/
sudo chmod +x /usr/local/bin/cleanup
```

Run from anywhere:

```bash
cleanup
```

## Controls

| Key | Action |
|-----|--------|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `enter` | Confirm mode |
| `q` / `ctrl+c` | Quit |
| any key | Exit done/error screen |

## Tasks

### Shallow (5)
1. Clear Trash — `rm -rf ~/.local/share/Trash/*`
2. Clear User Cache — `rm -rf ~/.cache/*`
3. Remove .tmp files — `find /tmp ~/.cache -name "*.tmp" -delete`
4. Vacuum journals (7 days) — `sudo journalctl --vacuum-time=7d`
5. APT autoclean — `sudo apt autoclean`

### Deep (12)
1. Save today's log — `sudo journalctl --since today > ~/cleanup-log-$(date +%F).txt`
2. APT update
3. APT upgrade -y
4. APT full-upgrade -y
5. APT autoremove -y
6. APT autoclean && apt clean
7. Clear /tmp and /var/tmp — `sudo rm -rf /tmp/* /var/tmp/*`
8. Clear user cache + trash — `rm -rf ~/.cache/* ~/.local/share/Trash/*`
9. Remove .tmp files — `find /tmp ~/.cache ~/.local -name "*.tmp" -delete 2>/dev/null`
10. Clean npm cache — `npm cache clean --force && rm -rf ~/.npm`
11. Prune pnpm store — `pnpm store prune && rm -rf ~/.pnpm-store`
12. Vacuum journals (3 days) — `sudo journalctl --vacuum-time=3d`

## Stack

- [Bubble Tea v2](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Bubbles v2](https://github.com/charmbracelet/bubbles) — spinner + progress
- [Lip Gloss v2](https://github.com/charmbracelet/lipgloss) — styling (Catppuccin Mocha palette)

## File layout

```
cleanup/
├── main.go      # entry point
├── model.go     # Bubble Tea model, screens, Update/View
├── tasks.go     # Task struct + ShallowTasks/DeepTasks
├── styles.go    # Catppuccin Mocha colors, lipgloss styles
├── go.mod
└── go.sum
```

## Notes

- Deep clean task 7 (`rm -rf /tmp/*`) clears non-hidden files in `/tmp`. Per-plan spec.
- On first task failure, remaining tasks are skipped and the error screen shows the cause + stderr.
