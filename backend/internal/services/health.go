package services

// Services is the application registry. The process constructs one in Run.
type Services struct{}

// New returns the application services.
func New() *Services {
	return &Services{}
}

// Health reports that the process is up. It does not check PostgreSQL or MinIO.
func (s *Services) Health() error {
	return nil
}
