# uncut CLI roadmap

All API methods from the original `methods/*.md` documentation are covered by
CLI commands in the current Go binary.

## Implemented coverage

| API area | Commands |
|---|---|
| Auth/config | `login`, `logout`, `config` |
| Wallet | `balance`, `deposit`, `wallet` |
| BIN catalog | `bins` |
| Cards | `new`, `cards`, `card`, `topup`, `withdraw`, `transactions`, `rename`, `phone`, `freeze`, `unfreeze`, `delete` |
| Operations | `operation`, `wait` |

## Future improvements

### Installer and releases

- Build release binaries for `darwin-arm64`, `darwin-amd64`, `linux-amd64`,
  and `linux-arm64`.
- Add a Homebrew formula that downloads prebuilt binaries, so users do not need
  Go installed.
- Add an install script that copies `uncut` into `~/.local/bin` or
  `/usr/local/bin`.

### UX aliases

Optional aliases to consider:

| Alias | Command |
|---|---|
| `uncut tx <card_id>` | `uncut transactions <card_id>` |
| `uncut ops <operation_id>` | `uncut operation <operation_id>` |
| `uncut remove <card_id>` | `uncut delete <card_id>` |

### Safer reveal mode

- Add `--copy number|cvv|expiry` to copy one field to clipboard on macOS.
- Add `--no-emoji` for terminals that do not render card glyphs.
- Add optional reveal confirmation for `uncut cards --reveal`.

### Output polish

- Add compact tables for `cards`, `wallet`, and `transactions`.
- Add `--csv` for export workflows.
- Add `--filter status=ACTIVE` for card lists.

### Tests

- Add a local mock HTTP server test suite through `UNCUT_BASE_URL`.
- Add golden output tests for all commands.
- Add non-interactive tests for destructive command confirmation.
