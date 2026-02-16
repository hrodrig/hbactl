# hbactl add — Sequence

Add a rule: backup, then append at end or insert after last rule for a user (`--after-user`). `--dry-run` only prints the line.

```mermaid
sequenceDiagram
    participant User
    participant hbactl
    participant PostgreSQL
    participant Filesystem

    User->>hbactl: hbactl add --type ... [--after-user X] [--dry-run]
    hbactl->>hbactl: validate type, method, addr (local vs host)

    alt --dry-run
        hbactl->>User: "would append" or "would insert after user X" + line
    else real add
        alt path not from --file
            hbactl->>PostgreSQL: connect
            hbactl->>PostgreSQL: SHOW hba_file
            PostgreSQL-->>hbactl: path
        end
        hbactl->>Filesystem: Backup(path) → .bak or .bak.<timestamp>
        Filesystem-->>hbactl: backup path
        hbactl->>User: Backup created at: ...
        alt --after-user set
            hbactl->>Filesystem: read file, find last line where user = afterUser
            hbactl->>Filesystem: insert new line after that line, write file
        else default
            hbactl->>Filesystem: append new line to file
        end
        hbactl->>User: Success. Run 'hbactl reload' to apply.
    end
```

[General](sequence-general.md) · [List](sequence-list.md) · [Remove](sequence-remove.md) · [Check](sequence-check.md) · [Reload](sequence-reload.md)
