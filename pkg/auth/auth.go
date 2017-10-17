package auth

type contextKey string

// Context keys to transport auth information.
const (
	ContextKeyUserID contextKey = "userID"
)
