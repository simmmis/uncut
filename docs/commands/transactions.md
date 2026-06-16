# `uncut transactions`

## Purpose

Show card transaction history, newest first.

## Synopsis

```sh
uncut transactions <card_id> [--page <n>|--all] [--raw|--json]
```

## API Mapping

```text
GET /cards/{card_id}/transactions?page=<n>
```

## Input

| Argument / Flag | Required | Description |
|---|---|---|
| `<card_id>` | yes | Card id |
| `--page <n>` | no | Page number, default `1` |
| `--all` | no | Fetch all pages |
| `--raw` | no | Print tab-separated rows without pagination text |
| `--json` | no | Print raw API JSON |

## Examples

Show latest card transactions:

```sh
uncut transactions card_demo_primary
```

Show page 2:

```sh
uncut transactions card_demo_primary --page 2
```

Use TSV for scripts:

```sh
uncut transactions card_demo_primary --raw
```

Ask the CLI for examples using current cards:

```sh
uncut help transactions
```

## Output

```text
OPENAI *CHATGPT
time: 2026-06-11T12:34:56+00:00
type: transaction.auth.approved
country: US
amount: $20.00
pre-auth: 20.00
posted: 0.00

source: provider
page 1, 10 per page
```

If more pages exist:

```text
next: uncut transactions <card_id> --page 2
```

Raw columns:

```text
time	type	merchant_name	merchant_country	original_amount	original_currency	pre_auth_amount	posted_amount
```

## Errors

- `not_found`: run `uncut cards`.
- Provider fallback may return `source: local`.
