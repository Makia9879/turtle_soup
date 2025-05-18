package constant

import "fmt"

var (
	CacheKeyActivityToken = "activity_token:"
)

func GetActivityToken(token string) string {
	return fmt.Sprintf("%s%s", CacheKeyActivityToken, token)
}
