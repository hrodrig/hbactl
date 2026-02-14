# hbactl reload — Sequence

Ask PostgreSQL to reload configuration so changes to `pg_hba.conf` take effect without restart.

```mermaid
sequenceDiagram
    participant User
    participant hbactl
    participant PostgreSQL

    User->>hbactl: hbactl reload
    hbactl->>PostgreSQL: connect (DATABASE_URL / --conn)
    PostgreSQL-->>hbactl: OK
    hbactl->>PostgreSQL: SELECT pg_reload_conf()
    PostgreSQL-->>hbactl: OK
    hbactl->>User: Success: configuration reloaded.
```

[General](sequence-general.md) · [List](sequence-list.md) · [Add](sequence-add.md) · [Check](sequence-check.md)
