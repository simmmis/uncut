# `uncut docs`

## Purpose

Print embedded markdown documentation from inside the `uncut` binary.

This command is for users and agents that installed `uncut` through Homebrew or
another package and do not have the source repository checkout locally.

Install `uncut` with one command:

```sh
brew install simmmis/tap/uncut
```

## Synopsis

```sh
uncut docs [--list|all|readme|methods|errors|<command>]
```

## Input

| Argument | Required | Description |
|---|---|---|
| `--list` | no | List available embedded docs topics |
| `all` | no | Print all embedded markdown docs |
| `readme` | no | Print the embedded README |
| `methods` | no | Print CLI-to-API method mapping |
| `errors` | no | Print error handling docs |
| `<command>` | no | Print one command manual, for example `new` or `topup` |

## Examples

List docs topics:

```sh
uncut docs --list
```

Read the main docs:

```sh
uncut docs readme
```

Read card creation docs:

```sh
uncut docs new
```

Read money movement docs:

```sh
uncut docs topup
uncut docs withdraw
```

Read raw output columns for scripts:

```sh
uncut docs cards
uncut docs wallet
uncut docs transactions
```

Print every embedded markdown file:

```sh
uncut docs all
```

## Output

Markdown text is printed to stdout.

When installed through Homebrew, the same markdown docs are also installed under:

```text
$(brew --prefix uncut)/share/doc/uncut
```

## Online Docs

- Repository: <https://github.com/simmmis/uncut>
- Releases: <https://github.com/simmmis/uncut/releases>
- Homebrew tap: <https://github.com/simmmis/homebrew-tap>

## Errors

- Unknown topic: run `uncut docs --list`.
