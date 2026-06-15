# uncut CLI

`uncut` is a standalone command-line client for the Uncutt Cards API.

It is built as a Go binary, but users do not need Go installed. Runtime
dependencies are intentionally zero: no Node.js, npm, curl, or external packages
are required.

## Quick Start

From this repository:

```sh
./dist/uncut --version
./dist/uncut login
./dist/uncut balance
./dist/uncut cards
```

If `uncut` is in your `PATH`:

```sh
uncut login
uncut balance
uncut cards
```

## Install Locally

Install into `~/.local/bin`:

```sh
mkdir -p ~/.local/bin
cp dist/uncut ~/.local/bin/uncut
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

Or install system-wide:

```sh
sudo cp dist/uncut /usr/local/bin/uncut
```

## Authentication

Login stores the API key and API endpoint locally:

```sh
uncut login
```

The values are saved to:

```text
~/.config/uncut/config.json
```

The file is written with `0600` permissions.

Lookup order:

1. `UNCUT_API_KEY` and `UNCUT_BASE_URL`
2. `~/.config/uncut/config.json`

## Common Workflow

Check wallet balance:

```sh
uncut balance
```

Show deposit addresses:

```sh
uncut deposit
```

Show available BINs and fees:

```sh
uncut bins
```

Create a zero-balance card:

```sh
uncut new <bin_id> --wait
```

If you do not know the BIN:

```sh
uncut new
```

This prints the BIN table plus ready copy-paste create commands.

List cards safely:

```sh
uncut cards
```

Example output:

```text
Facebook Ads
id: card_demo_ads
💳 **** **** **** 3523
EXP:**/**  CVV:***  Active
Balance: $23.55
```

Reveal one card:

```sh
uncut card <card_id> --reveal
```

Example output:

```text
Facebook Ads
id: card_demo_ads
💳 DEMO-CARD-NUMBER
EXP:MM/YY  CVV:DEMO  Active
Balance: $23.55
```

Top up a card:

```sh
uncut topup <card_id> --amount 50 --wait
```

Show card transactions:

```sh
uncut transactions <card_id>
```

Get command-specific examples:

```sh
uncut help topup
uncut help card
uncut help new
```

When you are logged in, help for card/BIN commands includes ready commands with
your current card ids and BIN ids.

## Commands

| Command | Purpose |
|---|---|
| [`login`](docs/commands/login.md) | Save API key and endpoint |
| [`logout`](docs/commands/logout.md) | Remove saved API key |
| [`config`](docs/commands/config.md) | Show auth config status |
| [`balance`](docs/commands/balance.md) | Show wallet balance |
| [`deposit`](docs/commands/deposit.md) | Show USDT deposit addresses |
| [`wallet`](docs/commands/wallet.md) | Show wallet transaction history |
| [`bins`](docs/commands/bins.md) | Show available card BINs |
| [`new`](docs/commands/new.md) | Create a new card |
| [`cards`](docs/commands/cards.md) | List cards |
| [`card`](docs/commands/card.md) | Show one card |
| [`topup`](docs/commands/topup.md) | Move wallet funds to a card |
| [`withdraw`](docs/commands/withdraw.md) | Move card funds to wallet |
| [`transactions`](docs/commands/transactions.md) | Show card transaction history |
| [`rename`](docs/commands/rename.md) | Rename a card |
| [`phone`](docs/commands/phone.md) | Update 3DS phone |
| [`freeze`](docs/commands/freeze.md) | Freeze a card |
| [`unfreeze`](docs/commands/unfreeze.md) | Unfreeze a card |
| [`delete`](docs/commands/delete.md) | Close a card and refund balance |
| [`operation`](docs/commands/operation.md) | Show async operation status |
| [`wait`](docs/commands/wait.md) | Poll async operation status |

## Built-In Help

Every command supports examples:

```sh
uncut help <command>
uncut <command> --help
```

Examples:

```sh
uncut help topup
uncut help transactions
uncut help new
```

## JSON Mode

Most API commands support:

```sh
uncut cards --json
```

`--json` is intended for scripts and agents. It prints API JSON without
decorative formatting.

## Error Handling

`uncut` prints structured, actionable errors:

```text
error: wallet balance is too low
message: Wallet balance is lower than required amount
code: insufficient_balance
http: 422
hint: run `uncut balance`; create with `--topup 0` or fund the wallet with `uncut deposit`
```

See [docs/ERRORS.md](docs/ERRORS.md).

## Build

From `cli_v1`:

```sh
GOCACHE="$PWD/.gocache" go build -o dist/uncut .
```

The resulting `dist/uncut` binary can be copied to another compatible machine
without installing Go.
