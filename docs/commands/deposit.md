# `uncut deposit`

## Purpose

Show USDT deposit addresses for funding the wallet.

## Synopsis

```sh
uncut deposit [--raw|--json]
```

## API Mapping

```text
GET /wallet/deposit-addresses
```

## Input

| Flag | Required | Description |
|---|---|---|
| `--raw` | no | Print tab-separated fields without a header |
| `--json` | no | Print raw API JSON |

## Examples

Show deposit addresses:

```sh
uncut deposit
```

Use JSON for scripts:

```sh
uncut deposit --json
```

Use TSV for shell scripts:

```sh
uncut deposit --raw
```

## Output

```text
USDT deposit addresses

ETH USDT
0x4bbeEB066eD09B7AEd07bF39EEe0460DFa261520

TRON USDT
TN3W4H6rK2ce4vX9YnFQHwKENnHjoxb3m9
```

Raw columns:

```text
chain	token	address
```

## Errors

- `deposit_address_unavailable`: retry later.
- `401 unauthorized`: run `uncut login`.
