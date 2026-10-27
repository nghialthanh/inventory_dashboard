FROM node:22-alpine AS ui
WORKDIR /src/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

FROM golang:1.25-alpine AS api
WORKDIR /src
ENV GOTOOLCHAIN=auto
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 go build -o /out/inventory-api .

FROM alpine:3.21
RUN apk add --no-cache ca-certificates wget
WORKDIR /app
COPY --from=api /out/inventory-api /app/inventory-api
COPY --from=ui /src/frontend/dist /app/static
ENV HTTP_ADDR=:8080
ENV STATIC_DIR=/app/static
EXPOSE 8080
USER nobody
ENTRYPOINT ["/app/inventory-api"]
