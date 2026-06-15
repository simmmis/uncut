# `uncut card`

## Purpose

Show one card by id.

## Synopsis

```sh
uncut card <card_id> [--reveal|--full] [--json]
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

## Errors

- `not_found`: run `uncut cards` and copy the current id.
- `provider_error`: retry later if using `--reveal`.
