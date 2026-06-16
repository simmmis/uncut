# `uncut delete`

## Purpose

Close a card permanently and refund remaining balance to wallet.

## Synopsis

```sh
uncut delete <card_id> [--yes] [--raw|--json]
```

## API Mapping

```text
DELETE /cards/{card_id}
```

## Input

| Argument / Flag | Required | Description |
|---|---|---|
| `<card_id>` | yes | Card id |
| `--yes` | no | Skip interactive confirmation |
| `--raw` | no | Print tab-separated refund fields |
| `--json` | no | Print JSON |

## Examples

Delete with confirmation prompt:

```sh
uncut delete card_demo_primary
```

Delete without prompt:

```sh
uncut delete card_demo_primary --yes
```

Use TSV:

```sh
uncut delete card_demo_primary --yes --raw
```

Ask the CLI for examples using current cards:

```sh
uncut help delete
```

## Confirmation

Without `--yes`, the command asks:

```text
delete card 01jxp...? type "delete" to confirm:
```

Only the exact word `delete` confirms.

## Output

```text
card deleted
returned: $25.00
wallet balance: 270.50 USDT
```

Raw columns:

```text
returned_amount	returned_currency	wallet_balance
```

## Errors

- `provider_error`: card is not deleted; retry later.
- `not_found`: card already deleted or wrong id.
