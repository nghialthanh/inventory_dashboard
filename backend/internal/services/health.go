// Package services declares the feature interfaces.
// Logic implements them. Adapters — the HTTP API, and later events or third-party clients — call these interfaces and do not import logic.
// This package has no feature behavior and does not touch the database.
package services

// Health is process liveness. It does not check PostgreSQL or MinIO.
type Health interface {
	Health() error
}
