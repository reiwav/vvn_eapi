package settable

import (
	"eapi/o/file"
	"eapi/o/ref_request"
	"eapi/o/request"
	"eapi/o/user"
	"eapi/reiway/tibero"
)

func CreateAllTable() {
	user.NewTable()
	file.NewTable()
	request.NewTable()
	ref_request.NewTable()
	return
	var usr = user.User{
		BaseModel: tibero.BaseModel{
			ID: 1,
		},
		Login:          "admin",
		Password:       "admin",
		FirstName:      "Admin",
		LastName:       "Admin",
		Email:          "admin@local",
		Activated:      true,
		LangKey:        "vi",
		CreatedBy:      "system",
		LastModifiedBy: "admin",
		MinioEndpoint:  "http://103.143.207.79:8082",
		MinioKey:       "admin",
		MinioSecret:    "wooribank@minio",
		MinioUseSSL:    false,
		MinioBucket:    "ekyc",
		MinioPrefix:    "woori",
	}
	usr.Update(false)
	var usr1 = user.User{
		BaseModel: tibero.BaseModel{
			ID: 2,
		},
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
	usr1.Update(false)
}
