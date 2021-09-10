package admin

import (
	"github.com/reiwav/x/rest"

	"github.com/gin-gonic/gin"
)

type AdminServer struct {
	*gin.RouterGroup
	rest.JsonRender
}

func NewAdminServer(parent *gin.RouterGroup, name string) {
	// var s = AdminServer{
	// 	RouterGroup: parent.Group(name),
	// }
	//s.Use(mid.AuthBasicJwt("", true))
	//user.NewUserServer(s.Group("user"))
	//s.GET("history", s.handleHistory)
	//s.GET("revenue", s.handleRevenue)
}
