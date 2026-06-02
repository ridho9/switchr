# switchr

<img width="1548" height="1010" alt="image" src="https://github.com/user-attachments/assets/49cec3a0-5d07-46ab-9131-6f24c0987ad9" />

A terminal TUI for browsing and attaching to [herdr](https://github.com/herdr/herdr) sessions.

## Install

```sh
make install
```

Builds the binary and copies it + the launcher script to `~/.local/bin/`.
Set `PREFIX` to change the destination:

```sh
make install PREFIX=/usr/local/bin
```

Or via `go install` (binary only):

```sh
go install github.com/ridho9/switchr@latest
```

## Usage

### Interactive mode (default)

```sh
switchr
```

Opens a full-screen TUI listing your herdr sessions. Navigate with arrow keys or vim keys, press enter to attach. When you detach from a session, the list reappears. Press q/esc to quit.

### Pipe / print mode

```sh
switchr --print          # explicit flag
switchr | jq             # auto-detected when stdout is not a TTY
```

Outputs the selected session as JSON:

```json
{"command":"open-herdr","session_name":"mysession"}
```

### Terminal startup integration

Copy `contrib/term-launcher` to your PATH (e.g. `~/.local/bin/`) and configure your terminal to use it as the startup command:

| Terminal  | Config                                    |
|-----------|-------------------------------------------|
| Ghostty   | `command = /home/user/.local/bin/term-launcher` |
| Kitty     | `launch -- /home/user/.local/bin/term-launcher` |
| WezTerm   | `default_prog = { "/home/user/.local/bin/term-launcher" }` |
| Alacritty | `shell = { program = "/home/user/.local/bin/term-launcher" }` |

The launcher skips switchr and starts a normal shell when `$HERDR_SESSION` is already set (inside an existing herdr session).

### Keybindings

| Key            | Action                |
|----------------|-----------------------|
| ↑ / k          | Move up               |
| ↓ / j          | Move down             |
| 1–9, 0         | Jump to session       |
| Enter          | Attach to session     |
| q / esc / ^C   | Quit                  |
| r              | Refresh sessions      |

### Incompatible herdr server

<img width="1473" height="965" alt="image" src="https://github.com/user-attachments/assets/0c1bdb1f-38d5-4940-a1c5-e4b4d4d25fe7" />

If switchr detects that the running herdr daemon is using an incompatible protocol, it shows a restart modal instead of the session list.

In that modal:

- `y` — restart the daemon and reload sessions
- `n` — dismiss the modal and continue
- `q` or `Ctrl-C` — quit switchr

## License

MIT
