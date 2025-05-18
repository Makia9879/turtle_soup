package constant

import "fmt"

const (
	CacheKeyActivityToken = "activity_token:"
	CacheKeySessionToken  = "session_token:"
)

func GetActivityToken(token string) string {
	return fmt.Sprintf("%s%s", CacheKeyActivityToken, token)
}

func GetSessionToken(token string) string {
	return fmt.Sprintf("%s%s", CacheKeySessionToken, token)
}
