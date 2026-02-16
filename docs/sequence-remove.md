# hbactl remove — Sequence

Remove one rule by **`--index`** or all matching rules by **`--user`** [**`--db`**] or **`--addr`**. Backup, then delete the line(s); **`--dry-run`** only prints what would be removed.

```mermaid
sequenceDiagram
    participant User
    participant hbactl
    participant PostgreSQL
    participant Filesystem

    User->>hbactl: hbactl remove --index N | --user X [--db Y] | --addr A [--dry-run]
    hbactl->>hbactl: validate: one of --index, --user, --addr (not mixed)

    alt path not from --file
        hbactl->>PostgreSQL: connect
        hbactl->>PostgreSQL: SHOW hba_file
        PostgreSQL-->>hbactl: path
    end

    hbactl->>Filesystem: ParseFileWithLineNumbers(path)
    Filesystem-->>hbactl: rules with LineNo, Index

    alt by --index
        hbactl->>hbactl: find rule where Index == N → single RuleWithLine
    else by --user [--db] or --addr
        hbactl->>hbactl: filter rules: MatchesUser(user, db) or MatchesAddress(addr) → list
    end

    alt no match(es)
        hbactl->>User: error: no rule(s) matching criteria
    else --dry-run
        hbactl->>User: "would remove N rule(s):" + each #index (line X): line
    else real remove
        hbactl->>Filesystem: Backup(path)
        Filesystem-->>hbactl: backup path
        hbactl->>User: Backup created at: ...
        hbactl->>Filesystem: RemoveLines(path, lineNumbers)
        hbactl->>User: Success. Run 'hbactl reload' to apply.
    end
```

[General](sequence-general.md) · [List](sequence-list.md) · [Add](sequence-add.md) · [Check](sequence-check.md) · [Reload](sequence-reload.md)
