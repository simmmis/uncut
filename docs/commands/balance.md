# `uncut balance`

## Purpose

Show the current wallet balance.

## Synopsis

```sh
uncut balance [--raw|--json]
```

## API Mapping

```text
GET /wallet
```

## Input

| Flag | Required | Description |
|---|---|---|
| `--raw` | no | Print only the numeric balance, exactly as a script-friendly number |
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

Script-friendly number:

```sh
uncut balance --raw
```

## Output

Text:

```text
Balance: 24.00 USDT
```

Raw:

```text
24
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
