package middlewares

import (
	"github.com/gin-gonic/gin"
)

var services = []string{
	"airtime",
	"bundle",
}

func Check() gin.HandlerFunc{
	return func (c*gin.Context) {
		
	}
}
