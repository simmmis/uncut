# `uncut bins`

## Purpose

List available card BINs, supported currencies, and fees.

## Synopsis

```sh
uncut bins [--raw|--json]
```

## API Mapping

```text
GET /card-bins
```

## Input

| Flag | Required | Description |
|---|---|---|
| `--raw` | no | Print tab-separated fields without a header |
| `--json` | no | Print raw API JSON |

## Examples

Show BINs and fees:

```sh
uncut bins
```

Use JSON for scripts:

```sh
uncut bins --json
```

Use TSV for shell scripts:

```sh
uncut bins --raw
```

Get create-card examples using current BINs:

```sh
uncut help new
```

## Output

```text
id                          name                        currencies  issue      topup  auth               withdraw  wallet
bin_demo_sg  SG 559268 ($0.5 auth fee)   USD         2.00 USDT  0.00%  0.00% + 0.50 USDT  0.00%     no
```

## Fields

| Column | Meaning |
|---|---|
| `id` | BIN id used by `uncut new` |
| `currencies` | Supported card currencies |
| `issue` | Fixed card issue fee |
| `topup` | Card top-up fee |
| `auth` | Purchase/authorization fee |
| `withdraw` | Card withdrawal fee |
| `wallet` | Apple Pay / Google Pay / 3DS support |

Raw columns:

```text
id	name	wallet_support	currencies	issue_fee	topup_percent	topup_fixed	auth_percent	auth_fixed	withdraw_percent	withdraw_fixed
```

## Errors

- `provider_unavailable`: retry later.
- `401 unauthorized`: run `uncut login`.
