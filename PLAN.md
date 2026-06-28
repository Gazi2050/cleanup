# CLEANUP CLI - FINAL PLAN

## Go + Bubble Tea | Ubuntu System Cleanup Tool

### Version: 1.0.0 | Status: FINALIZED

---

## 1. PROJECT OVERVIEW

- **Name:** `cleanup`
- **Language:** Go 1.21+
- **UI Framework:** Bubble Tea + Lip Gloss
- **Purpose:** Beautiful, interactive TUI for Ubuntu system cleanup
- **Modes:** 2 (Shallow Clean, Deep Clean)
- **Install:** Global binary via `/usr/local/bin/cleanup`

---

## 2. TWO CLEANUP MODES

### MODE 1 — Shallow Clean (Daily Usage)

> Fast, safe, everyday cleanup. Non-destructive. No sudo needed for most tasks.
> Estimated Time: 10-20 seconds

| #   | Task                     | Command                                   |
| --- | ------------------------ | ----------------------------------------- |
| 1   | Clear Trash              | rm -rf ~/.local/share/Trash/\*            |
| 2   | Clear User Cache         | rm -rf ~/.cache/\*                        |
| 3   | Remove .tmp files        | find /tmp ~/.cache -name "\*.tmp" -delete |
| 4   | Vacuum journals (7 days) | sudo journalctl --vacuum-time=7d          |
| 5   | APT autoclean            | sudo apt autoclean                        |

---

### MODE 2 — Deep Clean (Weekly/Monthly)

> Full system cleanup. Aggressive. Requires sudo. Clears everything.
> Estimated Time: 60-120 seconds

| #   | Task                     | Command                                                        |
| --- | ------------------------ | -------------------------------------------------------------- |
| 1   | Save today's log         | sudo journalctl --since today > ~/cleanup*log*$(date +%F).txt  |
| 2   | APT update               | sudo apt update                                                |
| 3   | APT upgrade              | sudo apt upgrade -y                                            |
| 4   | APT full-upgrade         | sudo apt full-upgrade -y                                       |
| 5   | APT autoremove           | sudo apt autoremove -y                                         |
| 6   | APT autoclean + clean    | sudo apt autoclean && sudo apt clean                           |
| 7   | Clear /tmp and /var/tmp  | sudo rm -rf /tmp/_ /var/tmp/_                                  |
| 8   | Clear user cache + trash | rm -rf ~/.cache/_ ~/.local/share/Trash/_                       |
| 9   | Remove .tmp files        | find /tmp ~/.cache ~/.local -name "\*.tmp" -delete 2>/dev/null |
| 10  | Clean npm cache          | npm cache clean --force && rm -rf ~/.npm                       |
| 11  | Prune pnpm store         | pnpm store prune && rm -rf ~/.pnpm-store                       |
| 12  | Vacuum journals (3 days) | sudo journalctl --vacuum-time=3d                               |

---

## 3. UI/UX DESIGN

### Color Palette (Eye-relaxing, Dark terminal friendly)

| Element      | Color         | Hex      |
| ------------ | ------------- | -------- |
| Header       | Soft Cyan     | #89DCEB  |
| Mode Badge   | Soft Purple   | #CBA6F7  |
| Running Task | Soft Yellow   | #F9E2AF  |
| Done Task    | Soft Green    | #A6E3A1  |
| Error Task   | Soft Red      | #F38BA8  |
| Pending Task | Muted Gray    | #6C7086  |
| Progress Bar | Cyan → Purple | gradient |
| Border       | Surface Gray  | #313244  |
| Background   | Transparent   | terminal |

> Color theme inspired by Catppuccin Mocha — soft, muted, not harsh on eyes

---

### Mode Selection Screen (Entry Point)

```
╭─────────────────────────────────────────╮
│                                         │
│   🧹  Ubuntu Cleanup CLI                │
│                                         │
│   Select cleanup mode:                  │
│                                         │
│   ▶  🌿  Shallow Clean                  │
│          Fast daily cleanup (~20s)      │
│                                         │
│      🔥  Deep Clean                     │
│          Full system cleanup (~90s)     │
│                                         │
│   ↑↓ navigate   enter select   q quit  │
╰─────────────────────────────────────────╯
```

---

### Progress Screen (During Execution)

```
╭─────────────────────────────────────────╮
│                                         │
│   🧹  Ubuntu Cleanup CLI  [DEEP CLEAN]  │
│                                         │
│   ✅  Save today's log                  │
│   ✅  APT update                        │
│   ⠸   APT upgrade          ← spinner   │
│   ·   APT full-upgrade     ← pending   │
│   ·   APT autoremove                   │
│   ·   APT autoclean + clean            │
│   ·   Clear /tmp /var/tmp              │
│   ·   Clear user cache                 │
│   ·   Remove .tmp files                │
│   ·   Clean npm cache                  │
│   ·   Prune pnpm store                 │
│   ·   Vacuum journals                  │
│                                         │
│   ████████████░░░░░░░░  5/12 tasks     │
│                                         │
╰─────────────────────────────────────────╯
```

---

### Success Screen (After Completion)

```
╭─────────────────────────────────────────╮
│                                         │
│   ✨  Cleanup Complete!                 │
│                                         │
│   ✅  Save today's log                  │
│   ✅  APT update                        │
│   ✅  APT upgrade                       │
│   ✅  APT full-upgrade                  │
│   ✅  APT autoremove                    │
│   ✅  APT autoclean + clean             │
│   ✅  Clear /tmp /var/tmp               │
│   ✅  Clear user cache                  │
│   ✅  Remove .tmp files                 │
│   ✅  Clean npm cache                   │
│   ✅  Prune pnpm store                  │
│   ✅  Vacuum journals                   │
│                                         │
│   ████████████████████  12/12 ✓        │
│                                         │
│   🎉  Your system is fresh and clean!  │
│       Time taken: 1m 23s               │
│                                         │
│   press any key to exit                 │
╰─────────────────────────────────────────╯
```

### Error Screen

```
╭─────────────────────────────────────────╮
│                                         │
│   ⚠️   Task Failed                      │
│                                         │
│   ✅  Save today's log                  │
│   ✅  APT update                        │
│   ❌  APT upgrade        ← failed here  │
│   ·   (remaining tasks skipped)         │
│                                         │
│   Error: exit status 1                  │
│   Tip: Make sure sudo is available      │
│                                         │
│   press any key to exit                 │
╰─────────────────────────────────────────╯
```

---

## 4. FILE STRUCTURE

```
cleanup-cli/
├── main.go              # Entry point — boots Bubble Tea program
├── model.go             # Bubble Tea model, Init/Update/View
├── tasks.go             # Task definitions (shallow + deep)
├── styles.go            # All Lip Gloss styles and colors
├── go.mod               # Module definition
├── go.sum               # Auto-generated checksums
└── README.md            # Setup and usage guide
```

---

## 5. CODE ARCHITECTURE

```
main.go
  └── tea.NewProgram(NewModel())

model.go
  ├── Model struct
  │     ├── mode       string        // "shallow" | "deep"
  │     ├── screen     string        // "select" | "running" | "done" | "error"
  │     ├── tasks      []Task        // task list for chosen mode
  │     ├── current    int           // index of running task
  │     ├── spinner    spinner.Model
  │     ├── progress   progress.Model
  │     ├── startTime  time.Time
  │     └── errMsg     string
  │
  ├── Init()   → start spinner tick
  ├── Update() → handle key input + task messages
  └── View()   → render correct screen

tasks.go
  ├── ShallowTasks() []Task   → 5 daily tasks
  └── DeepTasks()    []Task   → 12 full tasks

  Task struct:
    ├── name    string   // display label
    ├── command string   // shell command
    └── status  string   // "pending" | "running" | "done" | "error"

styles.go
  ├── All lipgloss.Style definitions
  ├── Color constants (Catppuccin Mocha palette)
  ├── Box/border styles
  └── Progress bar gradient
```

---

## 6. SCREENS & FLOW

```
START
  │
  ▼
[Select Screen]
  │  ↑↓ to choose mode
  │  enter to confirm
  ├──── Shallow Clean ──────────────────────────────┐
  │                                                 │
  └──── Deep Clean ────────────────────────────────┐│
                                                   ││
                                                   ▼▼
                                          [Running Screen]
                                            │  spinner animates
                                            │  progress bar fills
                                            │  tasks tick ✅ one by one
                                            │
                                    ┌───────┴────────┐
                                    │                │
                                    ▼                ▼
                              [Done Screen]    [Error Screen]
                               time taken       error message
                               press any key    press any key
                                    │                │
                                    └───────┬────────┘
                                            ▼
                                          EXIT
```

---

## 7. DEPENDENCIES

```
go.mod:
  module cleanup-cli
  go 1.21

  require:
    github.com/charmbracelet/bubbletea  v0.27.0   // TUI framework
    github.com/charmbracelet/bubbles    v0.18.0   // spinner, progress components
    github.com/charmbracelet/lipgloss   v0.12.1   // styling
```

---

## 8. BUILD & INSTALL STEPS

```bash
# Step 1 — Install Go (if not installed)
sudo apt install golang-go
go version   # must be 1.21+

# Step 2 — Create project
mkdir -p ~/cleanup-cli
cd ~/cleanup-cli

# Step 3 — Add files (main.go, model.go, tasks.go, styles.go, go.mod)

# Step 4 — Download dependencies
go mod tidy

# Step 5 — Build (stripped for small size)
go build -ldflags="-s -w" -o cleanup .

# Step 6 — Test locally
./cleanup

# Step 7 — Install globally
sudo mv cleanup /usr/local/bin/
sudo chmod +x /usr/local/bin/cleanup

# Step 8 — Run from anywhere
cleanup
```

---

## 9. KEYBOARD CONTROLS

| Key        | Action                    |
| ---------- | ------------------------- |
| ↑ / k      | Move selection up         |
| ↓ / j      | Move selection down       |
| Enter      | Confirm selection         |
| q / Ctrl+C | Quit anytime              |
| Any key    | Exit on done/error screen |

---

## 10. COMPLETION CHECKLIST

### Phase 1 — Plan ✅

- [x] Define two modes (shallow + deep)
- [x] List all tasks per mode
- [x] Design UI screens
- [x] Plan file structure
- [x] Plan architecture
- [x] Define color palette
- [x] Define keyboard controls
- [x] Define build steps

### Phase 2 — Implementation ⏳

- [ ] styles.go (colors, borders, styles)
- [ ] tasks.go (shallow + deep task lists)
- [ ] model.go (bubble tea model + screens)
- [ ] main.go (entry point)
- [ ] go.mod (dependencies)
- [ ] README.md (documentation)

### Phase 3 — Build & Test ⏳

- [ ] go mod tidy
- [ ] go build
- [ ] Test shallow mode
- [ ] Test deep mode
- [ ] Test error handling
- [ ] Test Ctrl+C behavior
- [ ] Install globally
- [ ] Test from any directory

---

## SUMMARY

| Property       | Value                          |
| -------------- | ------------------------------ |
| Language       | Go 1.21+                       |
| TUI Framework  | Bubble Tea                     |
| Styling        | Lip Gloss (Catppuccin Mocha)   |
| Modes          | Shallow (5 tasks), Deep (12)   |
| Progress Bar   | Yes (gradient cyan → purple)   |
| Spinner        | Yes (animated on running task) |
| Error Handling | Yes (with message + tip)       |
| Timer          | Yes (shows total time taken)   |
| Global Install | Yes (/usr/local/bin/cleanup)   |
| Binary Size    | ~6-8MB (stripped)              |
| Files          | 6 (main, model, tasks, styles, |
|                | go.mod, README)                |
