package user

import (
	"fmt"

	"github.com/reiwav/x/rest"
)

func GetByLogin(username, password string) (usr *User, err error) {
	usr, err = GetByUserPass(username, password)
	if err != nil || usr == nil {
		return nil, rest.BadRequest("error.validation")
	}
	if !usr.Activated {
		return nil, rest.BadRequest("account is not active")
	}
	var pwdEncryt = decrypt(usr.Password.String())
	fmt.Println(pwdEncryt)
	if pwdEncryt != password {
		return nil, rest.BadRequest("error.validation")
	}
	// if err := auth.ComparePassword(string(usr.Password), password); err != nil {
	// 	fmt.Println("=== VAO LOI3: ", err.Error(), password)
	// 	return nil, rest.BadRequestValid(rest.BadRequest("error.validation"))
	// }
	usr.Password = ""
	return usr, nil
}
