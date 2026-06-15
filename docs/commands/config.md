# `uncut config`

## Purpose

Show auth configuration status without exposing the full API key or endpoint.

## Synopsis

```sh
uncut config
```

## Input

No arguments.

## Examples

Show active auth source:

```sh
uncut config
```

## Output

With a saved key:

```text
api key: <masked-api-key>
api key source: /Users/example/.config/uncut/config.json
endpoint: configured
endpoint source: /Users/example/.config/uncut/config.json
config: /Users/example/.config/uncut/config.json
```

Without a key or endpoint:

```text
api key: not configured
endpoint: not configured
config: /Users/example/.config/uncut/config.json
```

## Environment

If `UNCUT_API_KEY` or `UNCUT_BASE_URL` are set, they take precedence over the
config file values.
