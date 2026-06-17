package middleware

import "context"

func newContextWithValue(key contextKey, value interface{}) context.Context {
	return context.WithValue(context.Background(), key, value)
}
