package circuitbreaker

import (
	"github.com/gin-gonic/gin"
	"errors"
	"net/http"
)

func fallbackHandler(c *gin.Context) {
    c.JSON(503, gin.H{
        "error": "Service unavailable, fallback activated",
    })
}

func NewCBMiddleware(cb *CircuitBreaker, fallback func(c *gin.Context)) gin.HandlerFunc {
    return func(c *gin.Context) {

        if c.Request.Method != http.MethodGet {
            c.Next()
            return
        }

        handled := false

        cb.Execute(
            func() error {
                // позволяем пройти запрос
                c.Next()

                // если handler вернул ошибку → ловим её
                if len(c.Errors) > 0 {
                    return errors.New(c.Errors.String())
                }
                return nil
            },
            func(c *gin.Context) {
                handled = true
                fallback(c)
            },
            c,
        )

        if handled {
            c.Abort()
        }
    }
}
