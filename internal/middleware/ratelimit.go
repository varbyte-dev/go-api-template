package middleware

import (
	"math"
	"strconv"
	"sync"
	"time"

	"go-api-template/internal/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type ipLimiterStore struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	r        rate.Limit
	burst    int
}

func newIPLimiterStore(r rate.Limit, burst int) *ipLimiterStore {
	s := &ipLimiterStore{
		visitors: make(map[string]*visitor),
		r:        r,
		burst:    burst,
	}
	go s.cleanup()
	return s
}

func (s *ipLimiterStore) get(ip string) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.visitors[ip]
	if !ok {
		v = &visitor{limiter: rate.NewLimiter(s.r, s.burst)}
		s.visitors[ip] = v
	}
	v.lastSeen = time.Now()
	return v.limiter
}

func (s *ipLimiterStore) cleanup() {
	const ttl = 10 * time.Minute
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		s.mu.Lock()
		for ip, v := range s.visitors {
			if time.Since(v.lastSeen) > ttl {
				delete(s.visitors, ip)
			}
		}
		s.mu.Unlock()
	}
}

func rateHeader(r rate.Limit) string {
	return strconv.FormatFloat(float64(r), 'f', -1, 64)
}

func retryAfterSeconds(delay time.Duration) int {
	return int(math.Ceil(delay.Seconds()))
}

func denyRequest(c *gin.Context, r rate.Limit, retryAfter int) {
	c.Header("X-RateLimit-Limit", rateHeader(r))
	c.Header("X-RateLimit-Remaining", "0")
	c.Header("Retry-After", strconv.Itoa(retryAfter))
	utils.TooManyRequests(c, "rate limit exceeded — retry after "+strconv.Itoa(retryAfter)+"s")
	c.Abort()
}

func allowRequest(c *gin.Context, l *rate.Limiter, r rate.Limit) {
	c.Header("X-RateLimit-Limit", rateHeader(r))
	c.Header("X-RateLimit-Remaining", strconv.Itoa(int(l.Tokens())))
}

// RateLimit returns a per-IP rate-limiting middleware using a token-bucket
// algorithm. Each unique client IP gets its own independent limiter so that
// a single misbehaving client does not affect others.
//
// Stale IP entries (idle for >10 min) are evicted automatically every
// 5 minutes to prevent unbounded memory growth.
//
// Parameters:
//   - r     : sustained request rate (events per second). Use rate.Limit(n)
//             for n req/s, or rate.Every(d) for one request per duration d.
//   - burst : maximum number of requests allowed in a single instant.
//
// Usage examples:
//
//	// 60 req/s per IP with a burst of 120 (general API routes)
//	protected.Use(middleware.RateLimit(rate.Limit(60), 120))
//
//	// 10 req/min per IP with a burst of 5 (auth — brute-force protection)
//	auth.Use(middleware.RateLimit(rate.Every(6*time.Second), 5))
func RateLimit(r rate.Limit, burst int) gin.HandlerFunc {
	store := newIPLimiterStore(r, burst)

	return func(c *gin.Context) {
		l := store.get(c.ClientIP())

		res := l.Reserve()
		if !res.OK() {
			denyRequest(c, r, 1)
			return
		}

		if delay := res.Delay(); delay > 0 {
			res.Cancel()
			denyRequest(c, r, retryAfterSeconds(delay))
			return
		}

		allowRequest(c, l, r)
		c.Next()
	}
}

// GlobalRateLimit returns a single shared rate-limiting middleware that
// applies to ALL requests regardless of their origin IP. This acts as a
// hard ceiling on total server throughput and is the first line of defence
// against traffic spikes.
//
// Because the bucket is shared, pair this with RateLimit for per-IP fairness.
//
// Usage example:
//
//	// Cap the whole server at 500 req/s with a burst of 1000
//	router.Use(middleware.GlobalRateLimit(rate.Limit(500), 1_000))
func GlobalRateLimit(r rate.Limit, burst int) gin.HandlerFunc {
	l := rate.NewLimiter(r, burst)

	return func(c *gin.Context) {
		res := l.Reserve()
		if !res.OK() {
			denyRequest(c, r, 1)
			return
		}

		if delay := res.Delay(); delay > 0 {
			res.Cancel()
			denyRequest(c, r, retryAfterSeconds(delay))
			return
		}

		allowRequest(c, l, r)
		c.Next()
	}
}
