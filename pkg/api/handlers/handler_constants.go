package handlers

import "time"

const (
	USER_LIST_CACHE_KEY = "userlist"
	POST_LIST_CACHE_KEY = "postlist"

	CACHE_DURATION_SHORT  = 30 * time.Second
	CACHE_DURATION_MEDIUM = 1 * time.Minute
	CACHE_DURATION_LONG   = 2 * time.Minute
)
