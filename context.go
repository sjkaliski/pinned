package pinned

import (
	"context"
)

// key is unexported and prevents collisions with
// context keys in other packages.
type key int

// contextKey is the context key for the Version.
const contextKey key = 0

// NewContext returns a new Context carrying a Version.
func NewContext(ctx context.Context, v *Version) context.Context {
	return context.WithValue(ctx, contextKey, v)
}

// FromContext returns the Version in the context.
func FromContext(ctx context.Context) *Version {
	v, _ := ctx.Value(contextKey).(*Version)
	return v
}
