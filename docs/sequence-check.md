# hbactl check — Sequence

Validate `pg_hba.conf` using the `pg_hba_file_rules` view. Requires a connection.

```mermaid
sequenceDiagram
    participant User
    participant hbactl
    participant PostgreSQL

    User->>hbactl: hbactl check
    hbactl->>PostgreSQL: connect (DATABASE_URL / --conn)
    PostgreSQL-->>hbactl: OK
    hbactl->>PostgreSQL: SELECT line_number, error FROM pg_hba_file_rules WHERE error IS NOT NULL ORDER BY line_number
    PostgreSQL-->>hbactl: rows (or empty)
    alt no errors
        hbactl->>User: OK: no syntax errors in pg_hba.conf
    else errors found
        hbactl->>User: Error: syntax errors in pg_hba.conf: line N: message
    end
```

[General](sequence-general.md) · [List](sequence-list.md) · [Add](sequence-add.md) · [Reload](sequence-reload.md)
