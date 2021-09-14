package mid

import (
	"eapi/common"
	"eapi/o/user"
	"fmt"
	"time"

	"github.com/reiwav/x/rest"
	"github.com/reiwav/x/web"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type MyCustomClaims struct {
	UserID string `json:"user_id,omitempty"`
	jwt.StandardClaims
}

type MyToken struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

// Example creating a token using a custom claims type.  The StandardClaim is embedded
// in the custom type to allow for easy encoding, parsing and validation of standard claims.
func CustomClaimsType(param *MyCustomClaims) (*MyToken, error) {
	mySigningKey := []byte(common.KEY_TOKEN)
	var timeNow = time.Now()
	// Create the Claims
	var timeExp = timeNow.Add(time.Hour * 24).Unix()
	param.UserID = param.Id
	param.ExpiresAt = timeExp
	param.Issuer = "token_user"
	param.IssuedAt = timeNow.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, param)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		return nil, err
	}
	return &MyToken{
		Token: ss,
	}, err
}

const (
	PARAM_TOKEN = "param_token"
	PARAM_USER  = "param_user"
	PARAM_CUS   = "param_cus"
	PARAM_CLAIM = "claim"
)

func AuthBasicJwt(tyRole string, isSetUser bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var errResponse = map[string]interface{}{
			"status": "error",
		}
		var tokenParam = web.GetToken(ctx.Request)
		if tokenParam == "" {
			errResponse["error"] = "token not found"
			ctx.JSON(401, errResponse)
			ctx.Abort()
			return
		}
		var claim, u, err = errorChecking(tokenParam, tyRole)
		if err != nil {
			errResponse["error"] = err.Error()
			ctx.JSON(401, errResponse)
			ctx.Abort()
			return
		}
		if isSetUser {
			ctx.Set(PARAM_USER, u)
		}
		ctx.Set(PARAM_TOKEN, claim)
		ctx.Next()
	}
}

func GetMyClaims(ctx *gin.Context) *MyCustomClaims {
	var data, isExist = ctx.Get(PARAM_TOKEN)
	if isExist {
		var res = data.(*MyCustomClaims)
		return res
	}
	return nil
}
func GetMyUser(ctx *gin.Context) *user.User {
	var data, isExist = ctx.Get(PARAM_USER)
	if isExist {
		var res = data.(*user.User)
		return res
	}
	return nil
}

func getTokenJwt(tokenParam string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenParam, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(common.KEY_TOKEN), nil
	})

}

// An example of parsing the error types using bitfield checks
func errorChecking(tokenParam, tyRole string) (*MyCustomClaims, interface{}, error) {
	// Token from another example.  This token is expired
	var token, err = getTokenJwt(tokenParam)
	if token.Valid {
		if claim, ok := token.Claims.(*MyCustomClaims); ok {
			fmt.Println("USER ID: ", claim.UserID)
			var dataUser, err = user.GetByID(claim.UserID)
			fmt.Println("USER : ", dataUser)
			if err != nil {
				return nil, nil, rest.Unauthorized("user " + err.Error())
			}

			var timeNow = time.Now().Unix()
			if timeNow >= claim.ExpiresAt {
				return nil, nil, rest.Unauthorized("token expried")
			}
			return claim, dataUser, nil
		}
		err = rest.Unauthorized("token not claim")
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return nil, nil, rest.Unauthorized("That's not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return nil, nil, rest.Unauthorized("Timing is everything")
		}
		return nil, nil, rest.Unauthorized("Couldn't handle this token:")
	}
	return nil, nil, rest.Unauthorized("Couldn't handle this token:")

}
