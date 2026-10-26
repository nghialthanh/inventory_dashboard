package trace

import "context"

type ctxKey struct{}

// Header is the HTTP header that carries the trace id.
const Header = "X-Trace-ID"

// WithID stores trace_id on the context for handler → logic → repository.
func WithID(ctx context.Context, id string) context.Context {
	if id == "" {
		return ctx
	}
	return context.WithValue(ctx, ctxKey{}, id)
}

// ID returns the trace_id from ctx, or empty when unset.
func ID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	id, _ := ctx.Value(ctxKey{}).(string)
	return id
}
