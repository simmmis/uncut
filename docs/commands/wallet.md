# `uncut wallet`

## Purpose

Show wallet transaction history, newest first.

## Synopsis

```sh
uncut wallet [--page <n>|--all] [--json]
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

## Errors

- Invalid page: `--page` must be a positive integer.
- `401 unauthorized`: run `uncut login`.
