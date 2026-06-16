# `uncut rename`

## Purpose

Rename a card's local label.

## Synopsis

```sh
uncut rename <card_id> <new_name> [--raw|--json]
```

Alternative:

```sh
uncut rename <card_id> --name <new_name>
```

## API Mapping

```text
PATCH /cards/{card_id}
```

## Input

| Argument / Flag | Required | API field | Description |
|---|---|---|---|
| `<card_id>` | yes | path | Card id |
| `<new_name>` | yes | `name` | New local card label |
| `--name <name>` | alternative | `name` | Long-form name |
| `--raw` | no | local | Print updated card as one TSV row |
| `--json` | no | local | Print JSON |

## Examples

Rename a card:

```sh
uncut rename card_demo_primary 'Facebook Ads'
```

Long-form name:

```sh
uncut rename card_demo_primary --name 'Google Ads'
```

Ask the CLI for examples using current cards:

```sh
uncut help rename
```

Use TSV:

```sh
uncut rename card_demo_primary 'Facebook Ads' --raw
```

## Output

Updated safe card output:

```text
Anthropic Billing
id: card_demo_ads
💳 **** **** **** 4242
EXP:**/**  CVV:***  Active
Balance: $25.00
```

Raw columns are the same as `uncut card --raw`:

```text
id	name	mask	card_number	expiration_date	cvv	currency	balance	status	expire_month	expire_year	phone_3ds	created_at
```

## Errors

- `validation_failed`: name missing or too long.
- `not_found`: run `uncut cards`.
