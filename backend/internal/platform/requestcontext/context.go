package requestcontext

import "context"

type key uint8

const (
	projectKey key = iota
	testTokenKey
)

func WithProjectKey(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, projectKey, value)
}

func ProjectKey(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(projectKey).(string)
	return value, ok
}

func WithTestToken(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, testTokenKey, value)
}

func TestToken(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(testTokenKey).(string)
	return value, ok
}
