# `uncut balance`

## Purpose

Show the current wallet balance.

## Synopsis

```sh
uncut balance [--json]
```

## API Mapping

```text
GET /wallet
```

## Input

| Flag | Required | Description |
|---|---|---|
| `--json` | no | Print raw API JSON |

## Examples

Human-readable balance:

```sh
uncut balance
```

Script-friendly JSON:

```sh
uncut balance --json
```

## Output

Text:

```text
Balance: 24.00 USDT
```

JSON:

```json
{
  "data": {
    "balance": 24,
    "currency": "USDT"
  }
}
```

## Errors

- Missing auth: run `uncut login`.
- `401 unauthorized`: key invalid or revoked.
