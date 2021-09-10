package user

// import (
// 	"eapi/common"
// 	mid "eapi/mid"
// 	"eapi/o/branch"
// 	"eapi/o/user"

// 	"github.com/reiwav/x/rest"

// 	"github.com/gin-gonic/gin"
// )

// type userServer struct {
// 	*gin.RouterGroup
// 	rest.JsonRender
// }
// type userBranch struct {
// 	*user.User
// 	DataBranch *branch.Branch `json:"data_branch"`
// }

// func NewUserServer(parent *gin.RouterGroup) {
// 	var s = userServer{
// 		RouterGroup: parent,
// 	}
// 	s.Use(mid.AuthBasicJwt("", false))
// 	s.GET("list", s.handleList)
// 	s.GET("get", s.handleGet)
// 	s.POST("create", s.handleCreate)
// 	s.POST("update", s.handleUpdate)
// 	s.DELETE("delete", s.handleDelete)
// }

// func (s *userServer) handleList(ctx *gin.Context) {
// 	var brID = ctx.Query("branch_id")
// 	q := ctx.Request.URL.Query()
// 	//var roles = web.GetArrString("role", ",", q)
// 	var claim = mid.GetMyClaims(ctx)
// 	var types = []string{}

// 	if claim.Type == common.TypeAdmin {
// 		var typeData = q.Get("type")
// 		if typeData == "user" {
// 			types = []string{common.TypeAdmin}
// 		} else {
// 			types = []string{common.TypeEmployee}
// 		}
// 	} else {
// 		types = []string{common.TypeEmployee}
// 	}

// 	var grps, err = user.GetByType(brID, types)
// 	rest.AssertNil(err)
// 	var branches, _ = branch.GetAll()
// 	var userBras = make([]userBranch, 0)
// 	for _, u := range grps {
// 		var br *branch.Branch
// 		if u.BranchID != "" {
// 			for _, bra := range branches {
// 				if u.BranchID == bra.ID {
// 					br = bra
// 					break
// 				}
// 			}
// 		}
// 		var userBra = userBranch{
// 			User:       u,
// 			DataBranch: br,
// 		}
// 		userBras = append(userBras, userBra)
// 	}
// 	s.SendDataError(ctx, userBras, err)
// }

// func (s *userServer) handleGet(ctx *gin.Context) {
// 	var id = ctx.Query("id")
// 	var gr, err = user.GetByID(id)
// 	if err == nil {
// 		if gr.BranchID != "" {
// 			var br, _ = branch.GetByID(gr.BranchID)
// 			var data = userBranch{
// 				User:       gr,
// 				DataBranch: br,
// 			}
// 			s.SendData(ctx, data)
// 			return
// 		}
// 	}
// 	s.SendData(ctx, gr)
// }

// func (s *userServer) handleDelete(ctx *gin.Context) {
// 	var id = ctx.Query("id")
// 	var u, err = user.GetByID(id)
// 	rest.AssertNil(err)
// 	err = u.MarkDelete()
// 	rest.AssertNil(err)
// 	s.Success(ctx)
// }

// func (s *userServer) handleCreate(ctx *gin.Context) {
// 	var u *user.User
// 	rest.AssertNil(ctx.BindJSON(&u))
// 	err := u.CreateCheckPhone()
// 	rest.AssertNil(err)
// 	s.SendData(ctx, u)
// }
// func (s *userServer) handleUpdate(ctx *gin.Context) {
// 	var id = ctx.Query("id")
// 	var u *user.User
// 	rest.AssertNil(ctx.BindJSON(&u))
// 	u.ID = id
// 	err := u.Update()
// 	rest.AssertNil(err)
// 	s.SendData(ctx, u)
// }
