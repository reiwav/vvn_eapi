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
	var res = []struct {
		RequestType string `json:"request_type"`
	}{}
	var err = request.SelectDistinct(nil, cols, &res)
	var uls = []Ul{}
	if err != nil {
		uls = append(uls, Ul{
			ID:   1,
			Name: "EMPTY",
		})
	} else {
		for i, val := range res {
			uls = append(uls, Ul{
				ID:   i,
				Name: string(val.RequestType),
			})
		}
	}
	s.SendString(ctx, uls)
}

func (s *RequestServer) handleFakeCodes(ctx *gin.Context) {
	var cols = []string{"fake_code"}
	var res = []struct {
		FakeCode string `json:"fake_code"`
	}{}
	var err = request.SelectDistinct(nil, cols, &res)
	var uls = []Ul{}
	if err != nil {
		uls = append(uls, Ul{
			ID:   1,
			Name: "EMPTY",
		})
	} else {
		for i, val := range res {
			uls = append(uls, Ul{
				ID:   i,
				Name: string(val.FakeCode),
			})
		}
	}
	s.SendString(ctx, uls)
}
