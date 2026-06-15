# `uncut logout`

## Purpose

Remove the saved local API key.

## Synopsis

```sh
uncut logout
```

## Input

No arguments.

## Examples

Remove local saved auth:

```sh
uncut logout
```

## Output

```text
logout success!
```

## Side Effects

Deletes:

```text
~/.config/uncut/config.json
```

If the file is already absent, the command still succeeds.

## Errors

- Filesystem permission errors return exit code `1`.
