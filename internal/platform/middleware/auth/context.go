package auth

import "context"

type contextKey string

const userIDKey contextKey = "userID"

func UserId(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(userIDKey).(int64)
	return id, ok
}
