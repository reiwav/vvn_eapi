package auth

import (
	"eapi/mid"
	"eapi/o/user"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

type LoginApp struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	RememberMe bool   `json:"rememberMe"`
}

type DataLogin struct {
	*mid.MyToken
	User interface{} `json:"user"`
}

func (l *LoginApp) LoginUser() (*DataLogin, error) {
	u, err := user.GetByLogin(l.Username, l.Password)
	if err != nil {
		return nil, err
	}
	// var t, errT = l.getToken(u.ID, u.Type)
	// if err != nil {
	// 	return nil, errT
	// }
	var token, _ = mid.CustomClaimsType(&mid.MyCustomClaims{
		UserID: fmt.Sprintf("%v", u.ID),
		StandardClaims: jwt.StandardClaims{
			Subject: l.Username,
			Id:      fmt.Sprintf("%v", u.ID),
		},
	})
	return &DataLogin{
		MyToken: token,
		User:    u,
	}, nil
}
