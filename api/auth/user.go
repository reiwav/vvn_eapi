package auth

import (
	"github.com/reiwav/x/rest"

	"eapi/mid"
	"eapi/o/user"
	"eapi/reiway/tibero"

	"github.com/gin-gonic/gin"
)

type AuthenServer struct {
	*gin.RouterGroup
	rest.JsonRender
}

func NewAuthenServer(parent *gin.RouterGroup, name string) {
	var s = AuthenServer{
		RouterGroup: parent.Group(name),
	}
	//s.Use(mid.AuthPublic())
	s.POST("/reset-password/init", s.handleResetPassInit) //1
	s.Use(mid.AuthBasicJwt("", true)).POST("/change-password", s.ChangePass)
	// s.POST("/user/register", s.handleRegister)                                    //2
	// s.Use(mid.AuthRefreshTokenJwt()).POST("/refresh_token", s.handleRefreshToken) //3
	// s.Use(mid.AuthBasicJwt("", true))
	// s.POST("/user/logout", s.handleLogout)        //4
	// s.POST("/user/reset_pass", s.handleResetPass) //5
}

func (s *AuthenServer) ChangePass(ctx *gin.Context) { //1
	var body = struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}{}
	ctx.BindJSON(&body)
	if body.NewPassword == "" {
		rest.AssertNil(rest.BadRequest("valid.newPassword"))
	}
	var cl = mid.GetMyClaims(ctx)

	var u, err = user.GetByLogin(cl.Subject, body.CurrentPassword)
	rest.AssertNil(err)
	u.Password = tibero.String(body.NewPassword)
	err = u.Update(true)
	rest.AssertNil(err)
	s.SendData(ctx, body)
}

func (s *AuthenServer) handleResetPassInit(ctx *gin.Context) { //1
	var body string
	ctx.Bind(&body)
	s.SendData(ctx, body)
}

// func (s *AuthenServer) handleRefreshToken(ctx *gin.Context) { //1
// 	var l = struct {
// 		RefreshToken string `json:"refresh_token"`
// 	}{}
// 	rest.AssertNil(ctx.BindJSON(&l))
// 	var claimToken, isExist = ctx.Get(mid.PARAM_CLAIM)
// 	fmt.Println("=========", claimToken)
// 	if isExist {
// 		var tk, err = mid.RefreshToken(l.RefreshToken, claimToken.(*mid.MyCustomClaims))
// 		rest.AssertNil(err)
// 		s.SendData(ctx, tk)
// 		return
// 	}
// 	s.SendData(ctx, rest.Unauthorized("token not found param set"))

// }

func (s *AuthenServer) handleRegister(ctx *gin.Context) { //2

}

// func (s *AuthenServer) handleLogout(ctx *gin.Context) { //3
// 	var claim = mid.GetMyClaims(ctx)
// 	err := token.UpdateRevoke(claim.UserID)
// 	rest.AssertNil(err)
// 	s.SendData(ctx, nil)
// }

// func (s *AuthenServer) handleResetPass(ctx *gin.Context) { //4
// 	var claim = mid.GetMyUser(ctx)
// 	var body = struct {
// 		PassOld string `json:"password_old"`
// 		PassNew string `json:"password_new"`
// 	}{}
// 	rest.AssertNil(ctx.BindJSON(&body))
// 	fmt.Println(body.PassOld, body.PassNew)
// 	err := cAuth.UserUpdatePass(claim, body.PassOld, body.PassNew)
// 	rest.AssertNil(err)
// 	s.SendData(ctx, claim)
// }
