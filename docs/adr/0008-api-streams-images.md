# The API streams primary images

The browser loads image bytes from the API. MinIO is reachable only from the API. Vehicle JSON reports whether a primary image exists, not where the object is stored.

## Considered options

- A presigned MinIO URL in the JSON. Go moves fewer bytes, and the browser then depends on MinIO being published, signed, and unexpired.
