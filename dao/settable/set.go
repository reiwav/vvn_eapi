package settable

import (
	"eapi/o/file"
	"eapi/o/ref_request"
	"eapi/o/request"
	"eapi/o/user"
)

func CreateAllTable() {
	user.NewTable()
	file.NewTable()
	request.NewTable()
	ref_request.NewTable()
	return
	var usr = user.User{
		Login:          "admin",
		Password:       "admin",
		FirstName:      "Admin",
		LastName:       "Admin",
		Email:          "admin@local",
		Activated:      true,
		LangKey:        "vi",
		CreatedBy:      "system",
		LastModifiedBy: "admin",
		MinioEndpoint:  "http://103.143.207.45:8082",
		MinioKey:       "admin",
		MinioSecret:    "wooribank@minio",
		MinioUseSSL:    false,
		MinioBucket:    "ekyc",
		MinioPrefix:    "woori",
	}
	usr.Create()
	var usr1 = user.User{
		Login:          "user",
		Password:       "user",
		FirstName:      "user",
		LastName:       "user",
		Email:          "user@local",
		Activated:      true,
		LangKey:        "en",
		CreatedBy:      "system",
		LastModifiedBy: "admin",
		MinioEndpoint:  "",
		MinioKey:       "user",
		MinioSecret:    "wooribank@minio",
		MinioUseSSL:    false,
		MinioBucket:    "",
		MinioPrefix:    "",
	}
	usr1.Create()
}
