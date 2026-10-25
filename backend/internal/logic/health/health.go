package health

// Health is process liveness. It does not check PostgreSQL or MinIO.
type Health struct{}

// New returns the liveness check.
func New() Health { return Health{} }

// Health reports that the process is up.
func (Health) Health() error { return nil }
