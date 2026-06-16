# `uncut new`

## Purpose

Create a new virtual card. This is an asynchronous operation.

## Synopsis

```sh
uncut new <bin_id> --topup <amount> [--name <name>] [--currency <code>] [--3ds --phone <phone>] [--wait] [--raw|--json]
```

Equivalent long form:

```sh
uncut new --bin <bin_id> --topup <amount> [--name <name>] [--currency <code>]
```

## API Mapping

```text
POST /cards
GET /cards
```

`GET /cards` is used first to prevent duplicate card names.

## Defaults

| Field | Default |
|---|---|
| `--name` | `card-YYYYMMDD-HHMM` |
| `--currency` | `USD` |
| `--topup` | required, no default |

## Input

| Argument / Flag | Required | API field | Description |
|---|---|---|---|
| `<bin_id>` | yes | `bin_id` | BIN id from `uncut bins` |
| `--bin <bin_id>` | alternative | `bin_id` | Long-form BIN input |
| `--name <name>` | no | `name` | Local card label, must be unique |
| `--currency <code>` | no | `currency` | Card currency |
| `--topup <amount>` | yes | `topup_amount` | Initial balance, must be `> 0` |
| `--amount <amount>` | alternative | `topup_amount` | Alias for `--topup` |
| `--3ds` | no | `enable_3ds` | Enable 3DS SMS confirmations |
| `--phone <phone>` | if `--3ds` | `phone` | E.164 phone, e.g. `+10000000000` |
| `--wait` | no | local | Poll operation until terminal status |
| `--raw` | no | local | Print TSV: operation fields without decoration |
| `--json` | no | local | Print JSON |

## Examples

Show available BINs and ready create commands:

```sh
uncut new
```

Create a card with defaults and a required initial top-up:

```sh
uncut new bin_demo_sg --topup 25 --wait
```

Create a named card:

```sh
uncut new bin_demo_sg --name 'Facebook Ads' --topup 25 --wait
```

Equivalent alias:

```sh
uncut new bin_demo_sg --name 'Google Ads' --currency USD --amount 25 --wait
```

Ask the CLI for current real BIN examples:

```sh
uncut help new
```

## Output

Without `--wait`:

```text
operation: op_demo_create
status: new
Card issue queued. Poll GET /api/v1/operations/{operation_id} until status is "completed", then use its "card_id".
next: uncut wait op_demo_create
```

With `--wait`, final operation status is printed.

Raw output without `--wait`:

```text
op_demo_create	new
```

Raw output with `--wait`:

```text
op_demo_create	card_issue	completed	-27	card_demo_ads	2026-06-11T10:00:00+00:00	2026-06-11T10:00:40+00:00	
```

Raw columns without `--wait`: `operation_id`, `status`.

Raw columns with `--wait`: `operation_id`, `type`, `status`, `amount`,
`card_id`, `created_at`, `updated_at`, `error_message`.

## Missing BIN Helper

If no BIN is supplied:

```sh
uncut new
```

The command prints available BINs, defaults, and copy-paste commands:

```text
copy-paste create commands:
  uncut new bin_demo_sg --name 'card-20260613-1420' --topup 25 --wait
```

## Duplicate Name Error

Names are checked before creation. Comparison is case-insensitive and normalizes
whitespace.

```text
new failed: card name must be unique; "Facebook Ads" already exists
existing card: card_demo_ads
try: uncut new bin_demo_sg --name 'Facebook Ads-2' --currency USD --topup 25 --wait
```

## Common Errors

- `invalid_bin`: run `uncut bins`.
- `unsupported_currency`: choose a currency listed by `uncut bins`.
- `insufficient_balance`: lower `--topup` or fund with `uncut deposit`.
- Missing `--topup`: pass a positive amount, for example `--topup 25`.
- `enable_3ds_unsupported`: choose a BIN with `wallet` = `yes`.
- `invalid_phone`: use E.164 format.
