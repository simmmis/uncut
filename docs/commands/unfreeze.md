# `uncut unfreeze`

## Purpose

Unfreeze a frozen card so it can be charged again.

## Synopsis

```sh
uncut unfreeze <card_id> [--json]
```

## API Mapping

```text
POST /cards/{card_id}/unfreeze
```

## Input

| Argument / Flag | Required | Description |
|---|---|---|
| `<card_id>` | yes | Card id; card must be frozen |
| `--json` | no | Print JSON |

## Examples

Unfreeze a card:

```sh
uncut unfreeze card_demo_primary
```

Ask the CLI for examples using current cards:

```sh
uncut help unfreeze
```

## Output

Updated card with status `Active`.

## Errors

- `card_not_frozen`: choose a frozen card.
- `provider_error`: retry later.
- `not_found`: run `uncut cards`.
