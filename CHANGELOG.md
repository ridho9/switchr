# Changelog

## v0.2.0

### Features

- **Incompatible server detection** — `switchr` checks `herdr status --json` while loading sessions and surfaces a restart modal when the daemon protocol is incompatible.
- **Restart modal** — restart, dismiss, or quit directly from the modal; a successful restart reloads sessions automatically and shows a restart notice.

### Maintenance

- **Server status refactor** — grouped server-status fields into a `serverStatus` struct for cleaner state handling.

## v0.1.0

Initial release.

### Features

- **Two-column TUI** — session list (left) with workspace→tab→pane tree (right)
- **Interactive attach** — enter to attach, auto-returns to session list on detach
- **Pipe mode** — `--print` flag or auto-detected non-TTY output for scripting
- **Keyboard navigation** — arrows, vim keys (j/k), number quick-select (1–0)
- **Help bar** — keybindings visible at bottom of screen
- **Status indicators** — ●/○ running/stopped icons, attached/detached labels
- **Manual refresh** — press `r` to reload sessions without restarting
- **Scrolling** — long session lists scroll with `···` overflow indicators
- **Terminal launcher** — `contrib/term-launcher` script for startup integration with
  Ghostty, Kitty, WezTerm, and Alacritty
- **Theme-adaptive** — terminal palette ANSI colors, adapts to any colorscheme
- **`--version` flag** — works for both `make install` and `go install`
- **Window title** — changes to `herdr: {name}` during attached sessions

### Performance

- Concurrent `lsof` attachment checks (~1.4s → ~75ms for 19 sessions)
