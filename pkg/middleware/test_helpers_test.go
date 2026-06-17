package middleware

import "context"

type testContextKey string

func newContextWithValue(key contextKey, value interface{}) context.Context {
	return context.WithValue(context.Background(), key, value)
}
