# 10: Attach a primary image

**What to build:** While a vehicle is in stock, the manager attaches one primary image, replaces it, or removes it. A vehicle with no image still shows a placeholder. The picture loads from the dashboard's own address. After stock-out, that picture still loads, and attaching, replacing, or removing it is rejected.

**Blocked by:** 06 Stock out a vehicle

**Status:** ready-for-human

- [x] While the vehicle is in stock, the manager can attach one primary image, replace it, and remove it
- [x] A vehicle with no primary image remains valid, and the list shows a placeholder
- [x] The picture loads from the dashboard's own address
- [x] The vehicle payload reports whether a primary image exists and does not include a storage location
- [x] Asking for a missing primary image is not found
- [x] The upload is a single image file, JPEG, PNG, or WebP, at most 5 MiB. Any other upload is rejected
- [x] Attaching, replacing, or removing a primary image after stock-out is rejected
- [x] A primary image attached before stock-out still loads after the vehicle is stocked out
- [x] This behavior is verified through the HTTP API against PostgreSQL. Image bytes are uploaded and fetched through the image routes, and the test does not open MinIO itself
