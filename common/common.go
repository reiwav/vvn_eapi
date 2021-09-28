package common

import "eapi/reiway/tibero"

type ConfigSystem struct {
	TokenPublic string          `json:"token_public"`
	DB          tibero.ConfigDB `json:"db"`
	Port        string          `json:"port"`
	FolderAdmin string          `json:"folder_admin"`
	TimeExpires int64           `json:"time_expires"`
	UseMinio    bool            `json:"use_minio"`
}

const (
	KEY_TOKEN = "toiyeuvietnam"
)

var ConfigSystemCache = &ConfigSystem{}
