# hbactl demo — list, add, remove (dry-run)
# Run from repo root: go build -o hbactl . && vhs docs/demo.tape

Output docs/demo.gif

Set Shell "bash"
Set FontSize 14
Set Width 900
Set Height 520
Set Padding 20
Set TypingSpeed 10ms

# Clean prompt for recording
Type "export PS1='$ '"
Enter
Sleep 300ms

# Version
Type "./hbactl --version"
Enter
Sleep 2s
Enter
Sleep 1s

# List rules from file (no DB needed)
Type "./hbactl list -f sample-pg_hba.conf"
Enter
Sleep 2s
Enter
Sleep 1s

# Group by user
Type "./hbactl list -f sample-pg_hba.conf --group-by user"
Enter
Sleep 3s
Enter
Sleep 1s

# Dry-run: append
Type "./hbactl add --dry-run -f sample-pg_hba.conf --type host --db all --user app --addr 192.168.1.10/32 --method scram-sha-256"
Enter
Sleep 1s
Enter
Sleep 1s

# Dry-run: insert after user
Type "./hbactl add --dry-run -f sample-pg_hba.conf --type host --db app_ops --user ops_user --addr 10.0.0.5/32 --method md5 --after-user ops_user"
Enter
Sleep 2s
Enter
Sleep 1s

# Dry-run: remove by user (all databases)
Type "./hbactl remove -f sample-pg_hba.conf --user app_user --dry-run"
Enter
Sleep 2s
Enter
Sleep 1s

# Dry-run: remove by user and database
Type "./hbactl remove -f sample-pg_hba.conf --user app_user --db app_planning --dry-run"
Enter
Sleep 2s
Enter
Sleep 1s
