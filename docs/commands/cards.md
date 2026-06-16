# `uncut cards`

## Purpose

List all cards.

## Synopsis

```sh
uncut cards [--reveal|--full] [--raw|--json]
```

## API Mapping

Safe mode:

```text
GET /cards
```

Reveal mode:

```text
GET /cards
GET /cards/{card_id}/details
```

In reveal mode, one details request is made per card.

## Input

| Flag | Required | Description |
|---|---|---|
| `--reveal` | no | Show PAN/CVV |
| `--full` | no | Alias for `--reveal` |
| `--raw` | no | Print tab-separated fields without a header |
| `--json` | no | Print JSON |

## Examples

List cards safely:

```sh
uncut cards
```

List cards with full PAN/CVV:

```sh
uncut cards --reveal
```

Use JSON for scripts:

```sh
uncut cards --json
```

Use TSV for shell scripts:

```sh
uncut cards --raw
```

## Output

Safe output:

```text
Facebook Ads
id: card_demo_ads
💳 **** **** **** 3523
EXP:**/**  CVV:***  Active
Balance: $23.55
```

Reveal output:

```text
Facebook Ads
id: card_demo_ads
💳 DEMO-CARD-NUMBER
EXP:MM/YY  CVV:DEMO  Active
Balance: $23.55
```

Raw columns:

```text
id	name	mask	card_number	expiration_date	cvv	currency	balance	status	expire_month	expire_year	phone_3ds	created_at
```

Without `--reveal`, `card_number`, `expiration_date`, and `cvv` are empty.

## Security

Use `--reveal` only when needed. Full card data is sensitive and may be
audit-logged by the API.

## Errors

- `provider_error` or `provider_unavailable` can happen in reveal mode.
- `401 unauthorized`: run `uncut login`.
