# `uncut phone`

## Purpose

Update the 3DS phone number for a card.

## Synopsis

```sh
uncut phone <card_id> --phone <e164_phone> [--json]
```

Positional phone is also accepted:

```sh
uncut phone <card_id> +10000000000
```

## API Mapping

```text
PUT /cards/{card_id}/3ds-phone
```

## Input

| Argument / Flag | Required | API field | Description |
|---|---|---|---|
| `<card_id>` | yes | path | Card id |
| `--phone <phone>` | yes | `phone` | E.164 phone number |
| `--json` | no | local | Print JSON |

## Examples

Set 3DS phone:

```sh
uncut phone card_demo_primary --phone +10000000000
```

Equivalent positional phone:

```sh
uncut phone card_demo_primary +10000000000
```

Ask the CLI for examples using current cards:

```sh
uncut help phone
```

## Output

Updated safe card output.

## Errors

- `invalid_phone`: use E.164, e.g. `+10000000000`.
- `provider_error`: retry later.
- `not_found`: run `uncut cards`.
