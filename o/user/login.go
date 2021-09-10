package user

import (
	"fmt"

	"github.com/reiwav/x/rest"
)

func GetByLogin(username, password string) (usr *User, err error) {
	usr, err = GetByUserPass(username, password)
	if err != nil || usr == nil {
		fmt.Println("=== VAO LOI2: ", err.Error())
		return nil, rest.BadRequest("error.validation")
	}

	// if err := auth.ComparePassword(string(usr.Password), password); err != nil {
	// 	fmt.Println("=== VAO LOI3: ", err.Error(), password)
	// 	return nil, rest.BadRequestValid(rest.BadRequest("error.validation"))
	// }
	usr.Password = ""
	return usr, nil
}
