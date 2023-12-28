package helper

const (
	DatabaseContext = "db"
	CacheContext    = "cache"
)

type HTTPError struct {
	Message string `json:"message"`
}
