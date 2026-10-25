# The Docker Hub MinIO image is no longer published. This is the last community
# release that still ships a binary. It is not published to the host.
FROM alpine:3.21
ARG TARGETARCH
ARG MINIO_RELEASE=RELEASE.2025-09-07T16-13-09Z
RUN apk add --no-cache ca-certificates curl \
	&& curl -fsSL -o /usr/local/bin/minio \
		"https://github.com/minio/minio/releases/download/${MINIO_RELEASE}/minio.linux-${TARGETARCH}.${MINIO_RELEASE}" \
	&& chmod +x /usr/local/bin/minio \
	&& apk del curl
EXPOSE 9000
ENTRYPOINT ["minio"]
CMD ["server", "/data"]
