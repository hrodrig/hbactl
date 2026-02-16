# hbactl — General flow

High-level view: user runs a command; optionally connects to PostgreSQL to discover or validate; reads or writes `pg_hba.conf` on the filesystem.

```mermaid
sequenceDiagram
    participant User
    participant hbactl
    participant PostgreSQL
    participant pg_hba.conf

    User->>hbactl: hbactl <command> [flags]
    hbactl->>hbactl: resolve path (--file or SHOW hba_file)

    alt needs connection (list without -f, add without -f, remove without -f, check, reload)
        hbactl->>PostgreSQL: connect
        PostgreSQL-->>hbactl: OK
    end

    alt list
        hbactl->>pg_hba.conf: read
        hbactl->>User: table (# index for remove)
    else add
        hbactl->>pg_hba.conf: backup, then append or insert
        hbactl->>User: success (run reload)
    else remove
        hbactl->>pg_hba.conf: read, backup, remove line by index, write
        hbactl->>User: success (run reload)
    else check
        hbactl->>PostgreSQL: pg_hba_file_rules
        hbactl->>User: OK or errors
    else reload
        hbactl->>PostgreSQL: pg_reload_conf()
        hbactl->>User: success
    end
```

[List](sequence-list.md) · [Add](sequence-add.md) · [Remove](sequence-remove.md) · [Check](sequence-check.md) · [Reload](sequence-reload.md)
