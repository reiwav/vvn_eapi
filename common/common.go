package common

import "eapi/reiway/tibero"

type ConfigSystem struct {
	TokenPublic string          `json:"token_public"`
	DB          tibero.ConfigDB `json:"db"`
	Port        string          `json:"port"`
	FolderAdmin string          `json:"folder_admin"`
}

const (
	KEY_TOKEN = "toiyeuvietnam"
)

var ConfigSystemCache = &ConfigSystem{}
