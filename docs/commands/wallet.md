# `uncut wallet`

## Purpose

Show wallet transaction history, newest first.

## Synopsis

```sh
uncut wallet [--page <n>|--all] [--raw|--json]
```

## API Mapping

```text
GET /wallet/transactions?page=<n>
```

## Input

| Flag | Required | Description |
|---|---|---|
| `--page <n>` | no | Page number, default `1` |
| `--all` | no | Fetch all pages until `has_more=false` |
| `--raw` | no | Print tab-separated rows without pagination text |
| `--json` | no | Print raw API JSON |

## Examples

Show latest wallet activity:

```sh
uncut wallet
```

Show the next page:

```sh
uncut wallet --page 2
```

Fetch every page:

```sh
uncut wallet --all
```

Use TSV for scripts:

```sh
uncut wallet --raw
```

## Output

```text
card_topup  completed  -120.00
id: wallet_demo_txn
fee: 0.00
card: card_demo_primary
comment: Card top-up *DEMO
created: 2026-06-12T18:40:23+00:00

page 1, 10 per page
```

If more pages exist:

```text
next: uncut wallet --page 2
```

Raw columns:

```text
id	type	status	amount	fee	card_id	comment	created_at
```

## Errors

- Invalid page: `--page` must be a positive integer.
- `401 unauthorized`: run `uncut login`.
