package domain

import "time"

const (
	AccountTokenExpireDuration  = time.Hour * 24 * 30
	AccountTokenRefreshDuration = time.Hour * 24 * 7
)
