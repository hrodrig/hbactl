# Security Policy

## Supported Versions

We release security updates for the current minor version. Older minors are not maintained.

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |
| < 0.1   | :x:                |

## Reporting a Vulnerability

If you find a security issue, please report it responsibly:

- **Preferred:** Open a [GitHub Security Advisory](https://github.com/hrodrig/hbactl/security/advisories/new) (private until we decide how to address it).
- **Alternative:** Email the maintainer (see profile or commit history) with a clear description and steps to reproduce.

We will acknowledge the report and respond as soon as possible. If the issue is accepted, we will work on a fix and coordinate disclosure. We do not have a bug bounty program.

**Note:** hbactl edits `pg_hba.conf` and may run with elevated permissions; only use it in trusted environments and review changes before applying them.
