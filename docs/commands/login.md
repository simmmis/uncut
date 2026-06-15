# `uncut login`

## Purpose

Save an API key and API endpoint for future commands.

## Synopsis

```sh
uncut login
```

## Input

Interactive prompt:

```text
Enter API key:
Enter API endpoint:
```

## Examples

Login interactively:

```sh
uncut login
```

## Output

```text
login success!
```

## Side Effects

Creates or overwrites:

```text
~/.config/uncut/config.json
```

The file is written with `0600` permissions.

The endpoint is private account configuration. It is saved locally only and
must not be committed to git.

## Errors

- Empty key: exits with code `2`.
- Empty or invalid endpoint: exits with code `2`.
- Config write failure: exits with code `1`.
