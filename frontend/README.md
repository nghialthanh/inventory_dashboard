# Inventory screen

React screen for the inventory dashboard. Compose serves the production build with the API on port 8080.

## Development

Start the API with Docker Compose from the repository root so PostgreSQL is running. Then, from this directory:

```bash
npm install
npm run dev
```

Vite serves the screen on its own origin and proxies `/api` and `/health` to `http://localhost:8080`. The screen calls those paths; it does not call a second origin and it does not ask the manager to sign in.
