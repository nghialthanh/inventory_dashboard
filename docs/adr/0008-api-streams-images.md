# The API streams primary images

Superseded by [ADR-0010](0010-no-object-storage.md). There is no image route.

The browser loads image bytes from the API. MinIO is reachable only from the API. Vehicle JSON reports whether a primary image exists, not where the object is stored.

## Considered options

- A presigned MinIO URL in the JSON. Go moves fewer bytes, and the browser then depends on MinIO being published, signed, and unexpired.
