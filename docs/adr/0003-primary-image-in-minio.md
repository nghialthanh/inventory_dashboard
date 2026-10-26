# One optional primary image in MinIO

Superseded by [ADR-0010](0010-no-object-storage.md). The dashboard does not store vehicle images.

The inventory list is easier to scan with a photo, so each vehicle may have one primary image. The bytes live in MinIO. Postgres stores the object key. A vehicle with no image is still valid, and the list shows a placeholder.

## Considered options

- No image. Matches the written brief, and the list is harder to scan.
- A pasted image URL. Skips object storage, and it is not an upload.
- A gallery of many photos. A larger product than this assessment needs.
