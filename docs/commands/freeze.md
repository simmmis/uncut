# `uncut freeze`

## Purpose

Freeze an active card so new charges are declined.

## Synopsis

```sh
uncut freeze <card_id> [--raw|--json]
```

## API Mapping

```text
POST /cards/{card_id}/freeze
```

## Input

| Argument / Flag | Required | Description |
|---|---|---|
| `<card_id>` | yes | Card id; card must be active |
| `--raw` | no | Print updated card as one TSV row |
| `--json` | no | Print JSON |

## Examples

Freeze a card:

```sh
uncut freeze card_demo_primary
```

Ask the CLI for examples using current cards:

```sh
uncut help freeze
```

Use TSV:

```sh
uncut freeze card_demo_primary --raw
```

## Output

Updated card with status `Frozen`.

Raw columns are the same as `uncut card --raw`.

## Errors

- `card_not_active`: choose an active card.
- `provider_error`: retry later.
- `not_found`: run `uncut cards`.
