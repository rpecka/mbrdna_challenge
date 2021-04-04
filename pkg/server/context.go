package server

import "context"

type requestContextKey string

const (
	authTokenKey = requestContextKey("mbcvAuthToken")
	vehicleIDKey = requestContextKey("vehicleID")
)

func makeContext(ctx context.Context, authToken, vehicleID string) context.Context {
	ctx = context.WithValue(ctx, authTokenKey, authToken)
	return context.WithValue(ctx, vehicleIDKey, vehicleID)
}

func authToken(ctx context.Context) string {
	if id := ctx.Value(authTokenKey); id != nil {
		return id.(string)
	}
	return ""
}

func vehicleID(ctx context.Context) string {
	if id := ctx.Value(vehicleIDKey); id != nil {
		return id.(string)
	}
	return ""
}
