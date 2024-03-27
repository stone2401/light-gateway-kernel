package pcore

import (
	"errors"
	"net/http"

	"golang.org/x/time/rate"
)

func NewRateLimiter(limit int) Handler {
	limiter := rate.NewLimiter(rate.Limit(limit), limit)

	return func(r *http.Request) (code int, err error) {
		if !limiter.Allow() {
			return http.StatusTooManyRequests, errors.New("too many request")
		}
		return 200, nil
	}
}
