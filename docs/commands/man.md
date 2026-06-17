# `uncut man`

## Purpose

Print a full standalone manual from inside the `uncut` binary.

This is the built-in fallback for environments where `man uncut` is unavailable
or the package did not install system man pages.

Install `uncut` with one Homebrew command:

```sh
brew install simmmis/tap/uncut
```

## Synopsis

```sh
uncut man [topic]
```

## Input

| Argument | Required | Description |
|---|---|---|
| `<topic>` | no | Optional command topic, for example `new`, `card`, or `withdraw` |

## Examples

Read the full built-in manual:

```sh
uncut man
```

Read a command topic:

```sh
uncut man new
uncut man topup
uncut man withdraw
```

Read markdown docs instead:

```sh
uncut docs new
```

Use the system man page when installed:

```sh
man uncut
```

Open the public repository:

```text
https://github.com/simmmis/uncut
```

## Output

Plain text is printed to stdout.

`uncut man` is intentionally self-contained and does not need local markdown
files or the source repository.

## Online Docs

- Repository: <https://github.com/simmmis/uncut>
- Releases: <https://github.com/simmmis/uncut/releases>
- Homebrew tap: <https://github.com/simmmis/homebrew-tap>

## Errors

- Unknown topic: run `uncut docs --list`.
