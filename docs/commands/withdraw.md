# `uncut withdraw`

## Purpose

Move funds from a card back to wallet. This is asynchronous.

## Synopsis

```sh
uncut withdraw <card_id> --amount <amount> [--wait] [--json]
```

Positional amount is also accepted:

```sh
uncut withdraw <card_id> 20
```

## API Mapping

```text
POST /cards/{card_id}/withdraw
```

## Input

| Argument / Flag | Required | API field | Description |
|---|---|---|---|
| `<card_id>` | yes | path | Card id |
| `--amount <amount>` | yes | `amount` | Amount in card currency, must be `> 0` |
| `--wait` | no | local | Poll operation until terminal status |
| `--json` | no | local | Print JSON |

## Examples

Withdraw 20 from a card and wait:

```sh
uncut withdraw card_demo_primary --amount 20 --wait
```

Equivalent positional amount:

```sh
uncut withdraw card_demo_primary 20 --wait
```

Ask the CLI for examples using current cards:

```sh
uncut help withdraw
```

## Output

```text
operation: op_demo_withdraw
status: new
Card withdrawal queued. Poll GET /api/v1/operations/{operation_id} until status is "completed".
next: uncut wait op_demo_withdraw
```

## Errors

- `insufficient_card_balance`: run `uncut card <card_id>`.
- `not_found`: run `uncut cards`.
- `exchange_rate_unavailable`: retry later.
