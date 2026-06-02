# README screenshot fixtures

This folder holds the screenshot-only mock data used to capture the README images.
Nothing here is wired into the main app.

## Files

- `screenshot_mock.go` — a copy of the temporary mock loader and fake session tree data.
  - It has `//go:build ignore`, so Go will not compile it unless you copy it elsewhere and remove that line.
- `incompatible-server-status.json` — the server/client version pair that was used for the restart-modal screenshot.

## Current fixture set

Session list:
- `switchr`
- `build-systems`
- `docs`
- `infra`

Modal screenshot values:
- server: `0.6.5`
- client: `0.6.6`

## Reusing later

If you need the screenshot again, the quickest path is usually:

1. Copy `screenshot_mock.go` to the repo root.
2. Remove the `//go:build ignore` line.
3. Re-add the temporary early returns in `herdr.go` if you want to bypass live `herdr` calls.
