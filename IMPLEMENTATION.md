# uncut CLI implementation notes

## Binary strategy

`uncut` is implemented as a small Go binary.

Why Go:

- builds a real executable;
- does not require Node.js;
- does not shell out to `curl` for API requests;
- uses only the Go standard library;
- keeps HTTP, JSON parsing, config, polling, and formatting inside one binary.

The build target is:

```text
cli_v1/dist/uncut
```

## Source layout

```text
cli_v1/
  go.mod
  main.go
  METHODS.md
  IMPLEMENTATION.md
  ROADMAP.md
  docs/
  man/uncut.1
  dist/uncut
```

The implementation is currently one `main.go` file. Once command behavior
stabilizes, split it into packages for API, config, output, and command parsing.

## Embedded documentation

The binary embeds the markdown docs and man page source through Go `embed`.

Installed users can read docs without the repository checkout:

```sh
uncut help
uncut help new
uncut docs --list
uncut docs topup
uncut docs errors
uncut man
```

Release archives must include:

```text
uncut
README.md
METHODS.md
docs/
man/uncut.1
```

The Homebrew formula installs `uncut` to `bin`, `man/uncut.1` to `man1`, and the
markdown docs to `share/doc/uncut`.

## Runtime dependencies

None beyond the operating system. Users do not need Go, Node.js, npm, or curl.

## Build

From `cli_v1`:

```sh
GOCACHE="$PWD/.gocache" go build -o dist/uncut .
```

## API endpoint

There is no compiled default endpoint. The endpoint is private configuration and
must come from `uncut login` or `UNCUT_BASE_URL`.

Placeholder:

```text
<api_endpoint>
```

For local/mock testing:

```sh
UNCUT_BASE_URL=http://127.0.0.1:8080 ./dist/uncut cards
```

## Auth storage

`uncut login` stores the API key and endpoint in:

```text
~/.config/uncut/config.json
```

The file is written with `0600` permissions.

Credential lookup order:

1. `UNCUT_API_KEY` and `UNCUT_BASE_URL`
2. `~/.config/uncut/config.json`

## Smoke checks

Without calling the live API:

```sh
./dist/uncut --version
./dist/uncut
HOME=/tmp/uncut-empty ./dist/uncut balance
```

Read-only live checks after login:

```sh
./dist/uncut balance
./dist/uncut deposit
./dist/uncut bins
./dist/uncut cards
```

Mutating commands should be tested against a mock API or a disposable account:

```sh
./dist/uncut new ...
./dist/uncut topup ...
./dist/uncut withdraw ...
./dist/uncut freeze ...
./dist/uncut unfreeze ...
./dist/uncut delete ...
```
