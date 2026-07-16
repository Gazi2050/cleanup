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

## Table of Contents

- [Why cleanup](#why-cleanup)
- [Quick Start](#quick-start)
- [Usage](#usage)
- [Modes](#modes)
- [Stack](#stack)

## Why cleanup

Typing cleanup commands manually is a chore. `rm -rf /tmp/*`, `sudo apt autoclean`, `journalctl --vacuum-time=3d`... it's easy to mess up, and one typo can break things. Running system cleanup manually is tedious and error-prone — one destructive command with a small mistake can cause real damage.

cleanup eliminates that risk. It bundles all of this into a simple, interactive tool. Pick your mode, confirm, and watch it finish. No typing. No mistakes. Just reliable, automated cleanup with real-time progress and intelligent sudo handling — it checks for passwordless sudo first, then prompts only if needed.

## Quick Start

> [!WARNING]
> **Debian-based Linux only** (Ubuntu, Debian, Mint, Pop!_OS, etc.). Sudo access required.

### Install

Download from [releases](https://github.com/Gazi2050/cleanup/releases):

```bash
sudo mv cleanup_<version>_linux_amd64 /usr/local/bin/cleanup && sudo chmod +x /usr/local/bin/cleanup
cleanup --version
```

### Uninstall

```bash
sudo rm /usr/local/bin/cleanup
```

## Build from Source

Requires Go 1.26+.

```bash
go build -ldflags="-s -w" -o cleanup ./cmd/cleanup
sudo mv cleanup /usr/local/bin/ && sudo chmod +x /usr/local/bin/cleanup
```

## Usage

Run `cleanup` and select a mode:

- **Shallow** — 5 tasks, ~20s (safe for daily use)
- **Deep** — 12 tasks, ~90s (full system cleanup, requires sudo)

Navigate with `↑↓` / `kj`, confirm with `enter`, quit with `q` or `ctrl+c`.

## Modes

**Shallow Clean** — 5 tasks, ~20s `daily use`
- Clear trash, cache, .tmp files
- Vacuum journals (3 days)
- APT autoclean

**Deep Clean** — 12 tasks, ~90s `full system`
- Updates: APT update, APT full-upgrade, global npm + pnpm packages
- Cleanups: APT autoremove + clean, npm cache, pnpm store, /tmp + /var/tmp, user caches, .tmp files, vacuum journals (preserves live sockets)

A run never aborts on the first failure. Tasks whose required tool is missing (`apt`, `npm`, `pnpm`, `journalctl`, `find`, etc.) are auto-detected and skipped, and the run continues, ending on a summary with a `r retry failed` option. apt commands wait up to 5 minutes for the dpkg lock (so they no longer collide with `unattended-upgrades`) and run fully non-interactively. apt/npm/pnpm tasks get a 5-minute budget; other tasks cap at 60s so a hung command never freezes the TUI.

## Stack

- [Bubble Tea v2](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Bubbles v2](https://github.com/charmbracelet/bubbles) — UI components
- [Lip Gloss v2](https://github.com/charmbracelet/lipgloss) — Styling

<p align="center">Made with ❤️ for Linux users</p>
