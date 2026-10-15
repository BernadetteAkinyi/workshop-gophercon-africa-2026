# How to run this project

This is the quick, platform-honest guide to getting the workshop app running.
The full introduction is in [`README.md`](README.md); this file is just "make it go".

## Why `make run` was failing

`make` was not installed on this machine, so PowerShell reported:

```
make : The term 'make' is not recognized as the name of a cmdlet, function,
script file, or operable program.
```

That is the whole problem: the command `make` didn't exist, so nothing in the
`Makefile` ever ran. It has now been installed with:

```powershell
winget install --id ezwinports.make --accept-source-agreements --accept-package-agreements
```

Restart your terminal (or open a new one) so it picks up the updated `PATH`, then
check it:

```powershell
make --version
# GNU Make 4.4.1
```

## Prerequisites

| Tool   | Check            | Install (Windows)                        |
|--------|------------------|------------------------------------------|
| Go     | `go version`     | the `.msi` from https://go.dev/dl        |
| Git    | `git --version`  | `winget install Git.Git`                 |
| `make` | `make --version` | `winget install --id ezwinports.make`    |

The workshop needs Go 1.27.1 or newer (`go.mod`), Git, and `make`.

## One-time setup

From the project root, in a fresh terminal:

```powershell
make setup
```

This builds the Go tools from `vendor/` (nothing is downloaded) and installs the
app as `workshop` in `~/go/bin`. It takes a couple of minutes.

## Run it

### Interactive CLI (`make run`)

```powershell
make run
```

**This is an interactive menu — it waits for you to type and press Enter, and
appears to "hang" until you do.** That is not an error. From the menu you can
migrate, seed, ingest sightings, list, view stats, and start the web server.

If you just want the database filled and the site up, use the shortcut:

```powershell
go run . db seed      # migrate + seed + ingest in one go
make web              # serve on http://localhost:8080
```

### Other useful targets

```powershell
make run        # interactive CLI (mock GBIF / Open-Meteo, works offline)
make run-live   # same, but against the real APIs (rate limited)
make web        # serve the site on :8080
make web-hot    # same, with templ hot reload
make templ      # regenerate templ after editing a .templ file
make sqlc       # regenerate sqlc after editing a query or migration
make test       # run the tests
make key        # save the workshop's Claude key for Claude Code
make load       # load test the home page (run while make web is up)
make update     # pull the instructor's latest changes
make push       # commit everything and push to your fork
```

## Windows gotchas

`make` on Windows runs recipes through `cmd.exe`, and a couple of targets in the
`Makefile` are written for macOS/Linux only:

- **`make help`** uses `awk` — it will fail on Windows. Read the comments in the
  `Makefile` instead, or this file.
- **`make kill-ports`** uses `lsof` and a POSIX shell loop — it will fail on
  Windows. Kill a stuck port manually:

  ```powershell
  # find the process on a port (e.g. 8080) and stop it
  Get-NetTCPConnection -LocalPort 8080 -State Listen |
    Select-Object -ExpandProperty OwningProcess |
    ForEach-Object { Stop-Process -Id $_ -Force }
  ```

Everything else (`run`, `web`, `web-hot`, `templ`, `sqlc`, `test`, `format`,
`vendor`, `setup`) is a plain `go`/`git` command and works fine.

## Running without `make`

Every recipe is just a `go` or `git` command, so you can run them directly if you
prefer:

```powershell
go run .              # == make run
go run . --live       # == make run-live
go run . web          # == make web
go run ./cmd/mockapi  # == make mock
go test ./...         # == make test
```

## Troubleshooting

- **`make: command not found` / "not recognized"** — install it (see above) and
  open a new terminal so `PATH` is refreshed.
- **Installed `make` but it's *still* not recognized** — the terminal inherited
  the environment from before the install. A **new tab** in an already-open
  terminal (or an IDE like GoLand/VS Code, which keeps a stale environment for
  its whole life) will not pick up the change. Fully quit and reopen the
  terminal *application* or IDE. To fix the current session without restarting,
  run:
  ```powershell
  $env:Path += ";C:\Users\HomePC\AppData\Local\Microsoft\WinGet\Packages\ezwinports.make_Microsoft.Winget.Source_8wekyb3d8bbwe\bin"
  make --version
  ```
  If a brand-new PowerShell window still fails, sign out and back in so
  `explorer.exe` refreshes its environment.
- **`make run` seems to hang** — it's waiting for menu input. Type a command and
  press Enter, or use `Ctrl+C` to quit.
- **Port 8080 (or 8082/8083) already in use** — kill the listener (see the
  Windows gotchas section) and try again.
- **`no such table` / empty site** — the database hasn't been created yet. Run
  `go run . db seed`.
- **Messy database** — restore the committed copy with `git checkout workshop.db`.
- **Tool build issues** — run `make tools-online` to install the same tools at
  the same versions from the internet.
