# The API serves the UI on one origin

This is an interview assessment, so Compose stays small. It runs Postgres and the API. The API serves `/api` and the built screen on port 8080. A multi-stage image build copies that screen into the API. The browser has one origin and does not need CORS.

## Considered options

- A separate nginx container that proxies `/api`. One more process, and it does not exercise the inventory rules.
- The API and the UI on separate ports, with CORS. A second origin to explain for the same screen.

Accepted in [01](../../.scratch/inventory-dashboard/issues/01-open-the-empty-dashboard.md).
