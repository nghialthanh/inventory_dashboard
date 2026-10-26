# No object storage

Vehicle images are out of scope. The written brief treats object storage as optional, and the dashboard can identify aging stock without a photo. PostgreSQL is the only store. This supersedes ADR-0003 and ADR-0008.

## Considered options

- Keep MinIO for an optional primary image. An extra service, credentials, and routes for a feature the brief does not require.
- Store image bytes in Postgres or on the API disk. Still an image product the assessment does not ask for.
