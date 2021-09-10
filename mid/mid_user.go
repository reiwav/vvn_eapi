package mid

// import (
// 	"eapi/o/user"

// 	"github.com/reiwav/x/rest"

// 	"github.com/gin-gonic/gin"
// )

// func GetUserByToken(ctx *gin.Context) (*user.User, error) {
// 	var tok, err = GetToken(ctx)
// 	if err != nil {
// 		return nil, rest.Unauthorized("access token not found")
// 	}
// 	var u, err1 = user.GetByID(tok.UserID)
// 	if err1 != nil {
// 		return nil, rest.ErrorOK("user not found")
// 	}
// 	return u, nil
// }
