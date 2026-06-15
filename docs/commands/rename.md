# `uncut rename`

## Purpose

Rename a card's local label.

## Synopsis

```sh
uncut rename <card_id> <new_name> [--json]
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

## Output

Updated safe card output:

```text
Anthropic Billing
id: card_demo_ads
💳 **** **** **** 4242
EXP:**/**  CVV:***  Active
Balance: $25.00
```

## Errors

- `validation_failed`: name missing or too long.
- `not_found`: run `uncut cards`.
