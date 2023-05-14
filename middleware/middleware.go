package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func RateLimiter(formatted string) gin.HandlerFunc {
	// For scaling purposes use Redis as store
	store := memory.NewStore()
	rate, _ := limiter.NewRateFromFormatted(formatted)
	lim := limiter.New(store, rate)
	return func(c *gin.Context) {
		clientLimit, _ := lim.Get(c, string(lim.GetIP(c.Request)))
		if clientLimit.Reached {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "Too many requests. Try again later.",
			})
		}
		c.Next()
	}
}
