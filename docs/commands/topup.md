# `uncut topup`

## Purpose

Move funds from wallet to a card. This is asynchronous.

## Synopsis

```sh
uncut topup <card_id> --amount <amount> [--wait] [--raw|--json]
```

Positional amount is also accepted:

```sh
uncut topup <card_id> 50
```

## API Mapping

```text
POST /cards/{card_id}/topup
```

## Input

| Argument / Flag | Required | API field | Description |
|---|---|---|---|
| `<card_id>` | yes | path | Card id |
| `--amount <amount>` | yes | `amount` | Amount in USDT, must be `> 0` |
| `--wait` | no | local | Poll operation until terminal status |
| `--raw` | no | local | Print TSV operation fields |
| `--json` | no | local | Print JSON |

## Examples

Top up a card by 60 USDT and wait:

```sh
uncut topup card_demo_primary --amount 60 --wait
```

Equivalent positional amount:

```sh
uncut topup card_demo_primary 60 --wait
```

Incorrect:

```sh
uncut topup card_demo_primary --60 --wait
```

Use `--amount 60`, not `--60`.

Ask the CLI for examples using current cards:

```sh
uncut help topup
```

## Output

```text
operation: op_demo_topup
status: new
Card top-up queued. Poll GET /api/v1/operations/{operation_id} until status is "completed".
next: uncut wait op_demo_topup
```

Raw output without `--wait`:

```text
op_demo_topup	new
```

Raw columns without `--wait`: `operation_id`, `status`.

With `--wait --raw`, columns are `operation_id`, `type`, `status`, `amount`,
`card_id`, `created_at`, `updated_at`, `error_message`.

## Errors

- `insufficient_balance`: run `uncut balance`.
- `not_found`: run `uncut cards`.
- `validation_failed`: amount missing or not positive.
