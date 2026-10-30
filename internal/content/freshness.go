package content

import "time"

func IsStale(now, fetchedAt time.Time, ttl time.Duration) bool {
	if ttl <= 0 {
		return false
	}
	return !fetchedAt.IsZero() && now.Sub(fetchedAt) >= ttl
}
