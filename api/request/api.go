package request

import (
	"eapi/mid"
	"eapi/o/request"

	"github.com/reiwav/x/rest"

	"github.com/gin-gonic/gin"
)

type RequestServer struct {
	*gin.RouterGroup
	rest.JsonRender
}

func NewRequestServer(parent *gin.RouterGroup, name string) {
	var s = RequestServer{
		RouterGroup: parent.Group(name),
	}
	s.Use(mid.AuthBasicJwt("", false))
	s.GET("/requestTypes", s.handleType)   //1
	s.GET("/fakeCodes", s.handleFakeCodes) //2
}

type Ul struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (s *RequestServer) handleType(ctx *gin.Context) {
	var cols = []string{"request_type"}
	var reqs, _ = request.SelectDistinct(nil, cols)
	var uls = []Ul{}
	if len(reqs) == 0 {
		uls = append(uls, Ul{
			ID:   1,
			Name: "EMPTY",
		})
	} else {
		for i, val := range reqs {
			uls = append(uls, Ul{
				ID:   i,
				Name: string(val),
			})
		}
	}
	s.SendString(ctx, uls)
}

func (s *RequestServer) handleFakeCodes(ctx *gin.Context) {
	var cols = []string{"fake_code"}
	var reqs, _ = request.SelectDistinct(nil, cols)
	var uls = []Ul{}
	if len(reqs) == 0 {
		uls = append(uls, Ul{
			ID:   1,
			Name: "EMPTY",
		})
	} else {
		for i, val := range reqs {
			uls = append(uls, Ul{
				ID:   i,
				Name: string(val),
			})
		}
	}
	s.SendString(ctx, uls)
}
