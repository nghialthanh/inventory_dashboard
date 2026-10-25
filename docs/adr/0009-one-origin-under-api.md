# The API serves the UI on one origin

This is an interview assessment, so Compose stays small. It runs Postgres, MinIO, and the API. The API serves `/api` and the built screen on port 8080. A multi-stage image build copies that screen into the API. MinIO is not published. The browser has one origin and does not need CORS.

## Considered options

- A separate nginx container that proxies `/api`. One more process, and it does not exercise the inventory rules.
- The API and the UI on separate ports, with CORS. A second origin to explain for the same screen.
- Publishing MinIO and handing the browser a direct URL. Rejected when image bytes were kept behind the API.
