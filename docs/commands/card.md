# `uncut card`

## Purpose

Show one card by id.

## Synopsis

```sh
uncut card <card_id> [--reveal|--full] [--raw|--json]
```

## API Mapping

Safe mode:

```text
GET /cards/{card_id}
```

Reveal mode:

```text
GET /cards/{card_id}/details
```

## Input

| Argument / Flag | Required | Description |
|---|---|---|
| `<card_id>` | yes | Card id from `uncut cards` |
| `--reveal` | no | Show PAN/CVV |
| `--full` | no | Alias for `--reveal` |
| `--raw` | no | Print one tab-separated row |
| `--json` | no | Print raw API JSON |

## Examples

Show one card safely:

```sh
uncut card card_demo_primary
```

Reveal full payment details:

```sh
uncut card card_demo_primary --reveal
```

Use TSV for scripts:

```sh
uncut card card_demo_primary --raw
```

Ask the CLI for examples using current cards:

```sh
uncut help card
```

## Output

```text
Facebook Ads
id: card_demo_ads
💳 **** **** **** 3523
EXP:**/**  CVV:***  Active
Balance: $23.55
```

Raw columns:

```text
id	name	mask	card_number	expiration_date	cvv	currency	balance	status	expire_month	expire_year	phone_3ds	created_at
```

## Errors

- `not_found`: run `uncut cards` and copy the current id.
- `provider_error`: retry later if using `--reveal`.
