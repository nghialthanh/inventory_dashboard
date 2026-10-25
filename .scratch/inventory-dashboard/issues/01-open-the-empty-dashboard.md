# 01: Open the empty dashboard

**What to build:** The manager opens one local address and sees the inventory dashboard shell. The screen and the API share that address. A health check responds. PostgreSQL and a private MinIO are running behind the app, and MinIO is not reachable from the browser.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

- [x] Opening the local dashboard address shows the inventory screen shell
- [x] The screen and the API are served together on port 8080
- [x] A health check responds successfully outside the versioned API
- [x] PostgreSQL is the database behind the app
- [x] MinIO is running for later primary images, is not published to the host, and its credentials stay on the Compose network
- [x] The manager is not asked to sign in
- [x] During development the screen can be served on its own and still call the API on the same paths
