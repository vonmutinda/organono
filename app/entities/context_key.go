package entities

type ContextKey string

const (
	ContextKeyIpAddress ContextKey = "ipAddress"
	ContextKeyRequestID ContextKey = "requestId"
	ContextKeyTokenInfo ContextKey = "tokenInfo"
	ContextKeyUserAgent ContextKey = "userAgent"
)
