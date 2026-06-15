# `uncut wait`

## Purpose

Poll an asynchronous operation until it is completed, failed, or timed out.

## Synopsis

```sh
uncut wait <operation_id> [--interval <seconds>] [--timeout <seconds>] [--json]
```

## API Mapping

```text
GET /operations/{operation_id}
```

## Input

| Argument / Flag | Required | Default | Description |
|---|---|---|---|
| `<operation_id>` | yes | - | Operation id |
| `--interval <seconds>` | no | `3` | Poll interval; positive integer |
| `--timeout <seconds>` | no | `120` | Max wait time; positive integer |
| `--json` | no | false | Print final status JSON |

## Examples

Wait for operation completion:

```sh
uncut wait op_demo_topup
```

Wait with a longer timeout:

```sh
uncut wait op_demo_topup --interval 5 --timeout 180
```

Typical `topup --wait` already calls this internally:

```sh
uncut topup card_demo_primary --amount 60 --wait
```

## Output

While waiting:

```text
status: new
status: pending
```

On completion:

```text
operation: op_demo_create
type: card_issue
status: completed
card: card_demo_ads
```

## Errors

- Invalid `--interval` or `--timeout`: exit code `2`.
- Terminal operation status `error`: exit code `1`.
- Timeout: exit code `1`.
