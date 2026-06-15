# `uncut operation`

## Purpose

Show one asynchronous operation status.

## Synopsis

```sh
uncut operation <operation_id> [--json]
```

## API Mapping

```text
GET /operations/{operation_id}
```

## Input

| Argument / Flag | Required | Description |
|---|---|---|
| `<operation_id>` | yes | Operation id from `new`, `topup`, or `withdraw` |
| `--json` | no | Print raw API JSON |

## Examples

Show operation status:

```sh
uncut operation op_demo_topup
```

Use JSON:

```sh
uncut operation op_demo_topup --json
```

## Output

```text
operation: op_demo_create
type: card_issue
status: completed
amount: -27.00
card: card_demo_ads
created: 2026-06-11T10:00:00+00:00
updated: 2026-06-11T10:00:40+00:00
```

## Terminal Statuses

| Status | Meaning |
|---|---|
| `completed` | Operation succeeded |
| `error` | Operation failed |

## Errors

- `not_found`: operation id is wrong or belongs to another account.
