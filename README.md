<p align="center">
<pre>
 ██████╗██╗     ███████╗ █████╗ ███╗   ██╗██╗   ██╗██████╗ 
██╔════╝██║     ██╔════╝██╔══██╗████╗  ██║██║   ██║██╔══██╗
██║     ██║     █████╗  ███████║██╔██╗ ██║██║   ██║██████╔╝
██║     ██║     ██╔══╝  ██╔══██║██║╚██╗██║██║   ██║██╔═══╝ 
╚██████╗███████╗███████╗██║  ██║██║ ╚████║╚██████╔╝██║     
 ╚═════╝╚══════╝╚══════╝╚═╝  ╚═╝╚═╝  ╚═══╝ ╚═════╝ ╚═╝     
</pre>
</p>

<p align="center"><strong>A Linux cleanup TUI for tidy systems</strong></p>

cleanup is an interactive terminal tool that keeps your Linux system clean. Pick a mode, confirm, and it runs a curated sequence of cleanup tasks with live progress, smart sudo handling, and failure-aware stops — all from a single TUI.

## Table of Contents

- [Installation](#installation)
- [Build from Source](#build-from-source)
- [Features](#features)
- [Modes](#modes)
- [Controls](#controls)
- [Command Cheatsheet](#command-cheatsheet)

## Installation

Download the appropriate binary from the [latest release](https://github.com/Gazi2050/cleanup/releases), then follow the commands below.

Use `<version>` as the release tag (for example `v0.0.1`).

**Prerequisites:** Linux (apt-based distro). Sudo access is required — passwordless sudo is recommended for a silent run, but cleanup will prompt for a password interactively if needed.

---

### Linux

**Install (amd64)**

```bash
sudo mv cleanup_<version>_linux_amd64 /usr/local/bin/cleanup && sudo chmod +x /usr/local/bin/cleanup
```

**Install (arm64)**

```bash
sudo mv cleanup_<version>_linux_arm64 /usr/local/bin/cleanup && sudo chmod +x /usr/local/bin/cleanup
```

**Checksum (recommended)**

```bash
sha256sum cleanup_<version>_linux_amd64
sha256sum cleanup_<version>_linux_arm64
```

Compare output with `checksums.txt` from the release.

**Verify install** — `cleanup --version`

**Uninstall**

```bash
sudo rm /usr/local/bin/cleanup
```

## Build from Source

Requires Go 1.26+.

```bash
go build -ldflags="-s -w" -o cleanup ./cmd/cleanup
```

Install globally:

```bash
sudo mv cleanup /usr/local/bin/ && sudo chmod +x /usr/local/bin/cleanup
```

Run from anywhere:

```bash
cleanup
```

## Features

### Interactive Mode Selection

Choose between Shallow Clean (fast daily) or Deep Clean (full system) via a keyboard-driven TUI.

> **Keys:**
>
> - `↑` / `k` — Move up
> - `↓` / `j` — Move down
> - `enter` — Confirm mode

### Smart Sudo Handling

Checks for passwordless sudo first. If available, cleanup runs silently. If not, a masked password prompt appears and is validated before any task runs.

> **Flow:**
>
> - Runs `sudo -n true` to detect passwordless access
> - If that fails → prompts for password (masked input)
> - Validates via `sudo -S -v` before starting tasks

### Live Progress

Each task runs sequentially with real-time status: pending, running (spinner), or done. A progress bar tracks overall completion.

> **Indicators:**
>
> - `·` — Pending
> - Spinner — Running
> - `✅` — Done
> - `❌` — Failed

### Failure-Aware

On the first task failure, cleanup stops, skips remaining tasks, and shows an error card with the cause and command output.

> **On error:**
>
> - Failing task marked with `❌ ← failed here`
> - Error message + stderr displayed
> - Press any key to exit

## Modes

| Mode    | Tasks | Time  | Sudo                          |
| ------- | ----- | ----- | ----------------------------- |
| Shallow | 5     | ~20s  | partial (journalctl, apt)     |
| Deep    | 11    | ~90s  | yes                           |

### Shallow (5)

1. Clear Trash — `rm -rf ~/.local/share/Trash/*`
2. Clear User Cache — `rm -rf ~/.cache/*`
3. Remove .tmp files — `find /tmp ~/.cache -name "*.tmp" -delete 2>/dev/null || true`
4. Vacuum journals (3 days) — `sudo journalctl --vacuum-time=3d`
5. APT autoclean — `sudo apt autoclean`

### Deep (11)

1. APT update — `sudo apt update`
2. APT upgrade — `sudo apt upgrade -y`
3. APT full-upgrade — `sudo apt full-upgrade -y`
4. APT autoremove — `sudo apt autoremove -y`
5. APT autoclean + clean — `sudo apt autoclean && sudo apt clean`
6. Clear /tmp and /var/tmp — `sudo rm -rf /tmp/* /var/tmp/*`
7. Clear user cache + trash — `rm -rf ~/.cache/* ~/.local/share/Trash/*`
8. Remove .tmp files — `find /tmp ~/.cache ~/.local -name "*.tmp" -delete 2>/dev/null || true`
9. Clean npm cache — `npm cache clean --force && rm -rf ~/.npm`
10. Prune pnpm store — `pnpm store prune && rm -rf ~/.pnpm-store`
11. Vacuum journals (3 days) — `sudo journalctl --vacuum-time=3d`

## Controls

| Key           | Action                          |
| ------------- | ------------------------------- |
| `↑` / `k`     | Move up                         |
| `↓` / `j`     | Move down                       |
| `enter`       | Confirm mode / submit password  |
| `q` / `ctrl+c`| Quit                            |
| any key       | Exit done/error screen          |

## Command Cheatsheet

| Command             | What it does                          |
| ------------------- | ------------------------------------- |
| `cleanup`           | Launch interactive TUI mode selector. |
| `cleanup --version` | Print the installed version.          |
| `cleanup -v`        | Print the installed version (short).  |
| `cleanup version`   | Print the installed version.          |

## Stack

- [Bubble Tea v2](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Bubbles v2](https://github.com/charmbracelet/bubbles) — spinner + progress + text input
- [Lip Gloss v2](https://github.com/charmbracelet/lipgloss) — styling (Catppuccin Mocha palette)

## File layout

```
cleanup/
├── cmd/
│   └── cleanup/
│       └── main.go        # entry point, version flag
├── internal/
│   ├── tasks/
│   │   └── tasks.go       # Task struct, ShallowTasks/DeepTasks, ModeName
│   └── ui/
│       ├── model.go       # Bubble Tea model, screens, Update/View
│       ├── layout.go      # Box/card primitives
│       ├── styles.go      # Catppuccin Mocha lipgloss styles
│       └── theme.go       # palette + icons
├── go.mod
├── go.sum
└── .github/workflows/build.yml
```

<p align="center">Made with ❤️ for Linux users</p>
