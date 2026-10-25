# 10: Attach a primary image

**What to build:** While a vehicle is in stock, the manager attaches one primary image, replaces it, or removes it. A vehicle with no image still shows a placeholder. The picture loads from the dashboard's own address. After stock-out, that picture still loads, and attaching, replacing, or removing it is rejected.

**Blocked by:** 06 Stock out a vehicle

**Status:** ready-for-agent

- [ ] While the vehicle is in stock, the manager can attach one primary image, replace it, and remove it
- [ ] A vehicle with no primary image remains valid, and the list shows a placeholder
- [ ] The picture loads from the dashboard's own address
- [ ] The vehicle payload reports whether a primary image exists and does not include a storage location
- [ ] Asking for a missing primary image is not found
- [ ] The upload is a single image file, JPEG, PNG, or WebP, at most 5 MiB. Any other upload is rejected
- [ ] Attaching, replacing, or removing a primary image after stock-out is rejected
- [ ] A primary image attached before stock-out still loads after the vehicle is stocked out
- [ ] This behavior is verified through the HTTP API against PostgreSQL. Image bytes are uploaded and fetched through the image routes, and the test does not open MinIO itself
