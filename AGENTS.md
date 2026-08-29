# Repository Guidance

## Work-branch changelog

When updating `docs/work-branch-changelog.md`, describe the current committed
branch delta, not historical branch intent.

- Compare the working branch with the current remote-tracking base using
  `origin/main...HEAD`; do not rely on a potentially stale local `main` ref.
- Derive the merge base, file counts, insertion/deletion totals, and file
  statuses from that exact range.
- Use `git diff --name-status --find-renames origin/main...HEAD` to classify
  the Files table. Call a change a rename only when that command reports an
  `R` status in the current range.
- Move changes already present at the merge base to “Changes Already Absorbed
  From `main`”; remove stale branch-specific claims rather than preserving them
  as historical context.
- Before finishing, verify the document with `git diff --check` and ensure its
  stated range and statistics match the live diff output.
