# Error Handling

`uncut` treats errors as part of the interface. Errors should explain:

1. what failed;
2. what the API returned, when available;
3. what to do next.

## Format

Typical API error:

```text
error: wallet balance is too low
message: Wallet balance is lower than required amount
code: insufficient_balance
http: 422
hint: run `uncut balance`; lower --topup/--amount or fund the wallet with `uncut deposit`
```

Typical CLI usage error:

```text
topup failed: --amount must be a positive number
```

Typical network error:

```text
network error: lookup <api_host>: no such host
hint: check internet connection, DNS, or UNCUT_BASE_URL
```

## Local Preflight Errors

### Duplicate card name

Before creating a card, `uncut new` fetches existing cards and checks the
requested name. Card names are compared case-insensitively after trimming and
normalizing whitespace.

Example:

```text
new failed: card name must be unique; "Facebook Ads" already exists
existing card: card_demo_ads
try: uncut new bin_demo_sg --name 'Facebook Ads-2' --currency USD --topup 25 --wait
```

### Missing BIN

`uncut new` requires a BIN. If omitted, it prints current BINs and copy-paste
commands:

```text
error: --bin is required

available bins:
...

defaults:
  name: card-20260613-1420
  currency: USD
  topup: required, must be > 0

copy-paste create commands:
  uncut new 01... --name 'card-20260613-1420' --topup 25 --wait
```

### Invalid wait options

`--interval` and `--timeout` must be positive integers. They are validated
before any mutating API request is sent.

```text
topup failed: --interval must be a positive integer
```

## API Error Hints

| API code | Meaning | Suggested action |
|---|---|---|
| `unauthorized` | API key missing, invalid, revoked, or account inactive | Run `uncut login` again or set `UNCUT_API_KEY` |
| `not_found` | Resource does not exist or belongs to another account | Run `uncut cards` and copy the current id |
| `validation_failed` | Request body or query is invalid | Check command flags and formats |
| `insufficient_balance` | Wallet balance is too low | Run `uncut balance`; lower `--topup`/`--amount`; fund via `uncut deposit` |
| `insufficient_card_balance` | Card balance is too low | Run `uncut card <card_id>` and choose a smaller amount |
| `invalid_phone` | Phone is not E.164 | Use a value like `+10000000000` |
| `invalid_bin` | BIN is unknown or inactive | Run `uncut bins` or `uncut new` and copy a current BIN |
| `enable_3ds_unsupported` | 3DS requested for a BIN that does not support it | Choose a BIN where `wallet` is `yes`, or remove `--3ds` |
| `unsupported_currency` | Currency is not supported by the BIN | Run `uncut bins` and choose a listed currency |
| `card_not_active` | Card must be active | Run `uncut unfreeze <card_id>` or choose an active card |
| `card_not_frozen` | Card must be frozen | Run `uncut freeze <card_id>` or choose a frozen card |
| `provider_error` | Provider rejected or failed the request | Retry later; for delete, the card was not deleted |
| `provider_unavailable` | Provider is unavailable | Retry later |
| `card_issue_unavailable` | Issuing is temporarily disabled | Retry later |
| `exchange_rate_unavailable` | Currency conversion rate unavailable | Retry later |
| `deposit_address_unavailable` | Deposit address provider failed | Retry later |
| HTTP `429` | Rate limit exceeded | Wait a minute and retry |

## Exit Codes

| Code | Meaning |
|---|---|
| `0` | Success |
| `1` | API, network, or runtime error |
| `2` | CLI usage or preflight error |
| `3` | Missing API key |
