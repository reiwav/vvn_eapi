package mid

import (
	//"basic/o/push_token"

	"eapi/common"

	"github.com/reiwav/x/rest"
	"github.com/reiwav/x/web"

	// "fmt"
	// "g/x/web"
	"github.com/reiwav/x/mlog"

	"github.com/gin-gonic/gin"
)

var logMiddle = mlog.NewTagLog("middle")

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logMiddle.Error(err)
				var errResponse = map[string]interface{}{
					"message": err.(error).Error(),
					"path":    c.Request.URL.Path,
				}
				if httpError, ok := err.(rest.IHttpError); ok {
					errResponse["status"] = httpError.StatusCode()
					c.JSON(httpError.StatusCode(), errResponse)
				} else {
					errResponse["status"] = 500
					c.JSON(500, errResponse)
				}
				c.Abort()
			}
		}()
		c.Next()
	}
}

func AddHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,PATCH,OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set(
			"Access-Control-Allow-Headers",
			"Content-Type, Content-Range, Content-Disposition, Authorization",
		)
		c.Writer.Header().Set(
			"Access-Control-Allow-Credentials",
			"true",
		)
		//remember
		if c.Request.Method == "OPTIONS" {
			c.Writer.WriteHeader(200)
			return
		}
		c.Next()
	}
}

func AuthPublic() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var errResponse = map[string]interface{}{
			"status": "error",
		}
		var tokenPub = web.GetTokenPublic(ctx.Request)
		if tokenPub == "" || tokenPub != common.ConfigSystemCache.TokenPublic {
			errResponse["error"] = "Not system"
			ctx.JSON(401, errResponse)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
