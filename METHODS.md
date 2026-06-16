# uncut CLI v1 commands

## Command style

`uncut` uses readable single-command navigation:

```sh
uncut <command> [target] [flags]
```

The CLI does not expose old API method names like `card-create` or nested
groups like `cards create`.

Every command has examples:

```sh
uncut help <command>
uncut <command> --help
```

For card and BIN commands, help tries to load current account data and prints
ready copy-paste commands with real ids.

## API coverage

Every method from the original `methods/*.md` API documentation has a CLI
command:

| CLI command | API method | Endpoint |
|---|---|---|
| `uncut balance` | `wallet-balance` | `GET /wallet` |
| `uncut deposit` | `wallet-deposit-addresses` | `GET /wallet/deposit-addresses` |
| `uncut wallet` | `wallet-transactions` | `GET /wallet/transactions` |
| `uncut bins` | `card-bins-list` | `GET /card-bins` |
| `uncut new` | `card-create` | `POST /cards` |
| `uncut cards` | `card-list` | `GET /cards` |
| `uncut card <card_id>` | `card-get` | `GET /cards/{card_id}` |
| `uncut card <card_id> --reveal` | `card-details` | `GET /cards/{card_id}/details` |
| `uncut topup <card_id>` | `card-topup` | `POST /cards/{card_id}/topup` |
| `uncut withdraw <card_id>` | `card-withdraw` | `POST /cards/{card_id}/withdraw` |
| `uncut transactions <card_id>` | `card-transactions` | `GET /cards/{card_id}/transactions` |
| `uncut rename <card_id>` | `card-rename` | `PATCH /cards/{card_id}` |
| `uncut phone <card_id>` | `card-update-3ds-phone` | `PUT /cards/{card_id}/3ds-phone` |
| `uncut freeze <card_id>` | `card-freeze` | `POST /cards/{card_id}/freeze` |
| `uncut unfreeze <card_id>` | `card-unfreeze` | `POST /cards/{card_id}/unfreeze` |
| `uncut delete <card_id>` | `card-delete` | `DELETE /cards/{card_id}` |
| `uncut operation <operation_id>` | `operation-status` | `GET /operations/{operation_id}` |
| `uncut wait <operation_id>` | `operation-status` | `GET /operations/{operation_id}` |

## Auth

```sh
uncut login
uncut config
uncut logout
```

`uncut login` stores the API key and API endpoint at:

```text
~/.config/uncut/config.json
```

The config file is written with `0600` permissions.

Credential lookup order:

1. `UNCUT_API_KEY` and `UNCUT_BASE_URL`
2. `~/.config/uncut/config.json`

## Wallet

```sh
uncut balance
uncut deposit
uncut wallet
uncut wallet --page 2
uncut wallet --all
```

Wallet history is newest first. If the API reports more pages, output includes:

```text
next: uncut wallet --page 2
```

## BINs and card creation

```sh
uncut bins
uncut new <bin_id> --topup 25
uncut new <bin_id> --topup 25 --wait
uncut new <bin_id> --amount 25 --wait
uncut new <bin_id> --name "Facebook Ads" --currency USD --topup 25
uncut new <bin_id> --name "Facebook Ads" --currency USD --topup 25 --3ds --phone +10000000000
```

Defaults:

| Field | Default |
|---|---|
| `--name` | `card-YYYYMMDD-HHMM` |
| `--currency` | `USD` |
| `--topup` | required, positive number |

`--bin <bin_id>` is still supported, but the shorter form is:

```sh
uncut new <bin_id> --topup 25
```

If the BIN is missing, `uncut new` fails and prints available BINs with prices
plus ready copy-paste commands:

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

Before creating a card, `uncut new` also checks existing cards and refuses
duplicate names:

```text
new failed: card name must be unique; "Facebook Ads" already exists
existing card: card_demo_ads
try: uncut new 01... --name 'Facebook Ads-2' --currency USD --topup 25 --wait
```

## Cards

Safe output:

```sh
uncut cards
uncut card <card_id>
```

Reveal full PAN/CVV:

```sh
uncut cards --reveal
uncut card <card_id> --reveal
```

Safe card output:

```text
Facebook Ads
id: card_demo_ads
💳 **** **** **** 3523
EXP:**/**  CVV:***  Active
Balance: $23.55
```

Revealed card output:

```text
Facebook Ads
id: card_demo_ads
💳 DEMO-CARD-NUMBER
EXP:MM/YY  CVV:DEMO  Active
Balance: $23.55
```

The current API docs do not include a daily card limit field, so v1 does not
print `Limit: ...`. Add it only if the API starts returning such a field.

## Card money movement

```sh
uncut topup <card_id> --amount 50
uncut topup <card_id> --amount 50 --wait
uncut topup <card_id> 50 --wait
uncut withdraw <card_id> --amount 20
uncut withdraw <card_id> --amount 20 --wait
```

Do not write amounts as flags:

```sh
# wrong
uncut topup <card_id> --50 --wait

# right
uncut topup <card_id> --amount 50 --wait
```

`new`, `topup`, and `withdraw` are asynchronous and return an operation id.

## Transactions

```sh
uncut transactions <card_id>
uncut transactions <card_id> --page 2
uncut transactions <card_id> --all
```

The API returns newest transactions first, so default `page=1` is the latest
page. If more pages exist, output includes:

```text
next: uncut transactions <card_id> --page 2
```

## Card management

```sh
uncut rename <card_id> "New name"
uncut phone <card_id> --phone +10000000000
uncut freeze <card_id>
uncut unfreeze <card_id>
uncut delete <card_id>
uncut delete <card_id> --yes
```

`delete` asks for confirmation unless `--yes` is passed.

## Operations

```sh
uncut operation <operation_id>
uncut wait <operation_id>
uncut wait <operation_id> --interval 5 --timeout 180
```

`wait` polls until `completed`, `error`, or timeout.

## JSON mode

API commands support:

```sh
--json
```

`--json` prints machine-readable JSON without decorative formatting.

## Raw mode

API commands support:

```sh
--raw
```

`--raw` prints shell-friendly values. `uncut balance --raw` prints a single
number exactly as returned by the API, for example `49.3`. List and object
commands print tab-separated rows without a header. Command-specific manuals in
`docs/commands/` define the columns.

## Command manuals

Detailed command manuals live in:

```text
docs/commands/
```

General error behavior is documented in:

```text
docs/ERRORS.md
```
