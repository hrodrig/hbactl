# hbactl list — Sequence

List rules from `pg_hba.conf` in a table. Path from `--file` or from PostgreSQL (`SHOW hba_file`).

```mermaid
sequenceDiagram
    participant User
    participant hbactl
    participant PostgreSQL
    participant Filesystem

    User->>hbactl: hbactl list [--sort col] [--group-by user]
    alt path from --file
        hbactl->>hbactl: use -f path (no connection)
    else no --file
        hbactl->>PostgreSQL: connect (DATABASE_URL / --conn)
        PostgreSQL-->>hbactl: OK
        hbactl->>PostgreSQL: SHOW hba_file
        PostgreSQL-->>hbactl: path
    end
    hbactl->>Filesystem: read file
    Filesystem-->>hbactl: content
    hbactl->>hbactl: parse rules
    opt --sort or --group-by
        hbactl->>hbactl: sort by column
    end
    opt --group-by user
        hbactl->>hbactl: output with "=== user: X ===" separators
    end
    hbactl->>User: File: path (N rules) + table
```

[General](sequence-general.md) · [Add](sequence-add.md) · [Check](sequence-check.md) · [Reload](sequence-reload.md)
