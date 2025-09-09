package config

import (
	"fmt"
	"path/filepath"
)

const rootDir = "/media/mahdi/Cloud/Happle"
const applicationDir = "com.helium.message"
const usersDir = "users"

const UserId = "018f3a8b-1b32-7290-b1d5-92716a445330"
const ChatID = "018f3a8b-1b32-7295-a2c7-87654b4d4567"
const MessageID = "01991bc4-faad-7b70-aedc-f20ea4146898"

const Parsa = "018f3a8b-1b32-7291-a1c8-29817a544561"

func GetPath(file string) string {
	return filepath.Join(rootDir, applicationDir, file)
}

func GetUserPath(phone string, file string) string {
	pp := filepath.Join(rootDir, applicationDir, usersDir, phone, file)
	fmt.Println(pp)
	return pp
}
