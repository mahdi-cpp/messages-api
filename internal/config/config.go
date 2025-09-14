package config

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/google/uuid"
)

const rootDir = "/app/iris/"
const applicationDir = "com.iris.messages"
const usersDir = "users"

var (
	Mahdi  uuid.UUID
	Parsa  uuid.UUID
	Ali    uuid.UUID
	Golnar uuid.UUID
)
var (
	Digikala uuid.UUID
	Varzesh3 uuid.UUID
)
var (
	ChatID    uuid.UUID
	MessageID uuid.UUID
)

func initUsers() {
	var err error

	Mahdi, err = uuid.Parse("018f3a8b-1b32-7290-b1d5-92716a445330")
	if err != nil {
		log.Fatalf("failed to parse Mahdi: %v", err)
	}

	Parsa, err = uuid.Parse("018f3a8b-1b32-7291-a1c8-29817a544561")
	if err != nil {
		log.Fatalf("failed to parse Mahdi: %v", err)
	}

	Ali, err = uuid.Parse("018f3a8b-1b32-729c-a1b2-9876a5b4c3d2")
	if err != nil {
		log.Fatalf("failed to parse Mahdi: %v", err)
	}

	Golnar, err = uuid.Parse("018f3a8b-1b32-729f-d4e5-918273645a2c")
	if err != nil {
		log.Fatalf("failed to parse Mahdi: %v", err)
	}

}

func initChats() {
	var err error

	Digikala, err = uuid.Parse("018f3a8b-1b32-7292-b2d9-1237a7b8c8d2")
	if err != nil {
		log.Fatalf("failed to parse Mahdi: %v", err)
	}
	Varzesh3, err = uuid.Parse("018f3a8b-1b32-7293-c1d4-8765f5d1e2f3")
	if err != nil {
		log.Fatalf("failed to parse Mahdi: %v", err)
	}

}

// The Init function is called before main() and is ideal for initialization
func Init() {
	var err error

	initUsers()

	ChatID, err = uuid.Parse("018f3a8b-1b32-7295-a2c7-87654b4d4567")
	if err != nil {
		log.Fatalf("failed to parse ChatID: %v", err)
	}

	MessageID, err = uuid.Parse("01991bc4-faad-7b70-aedc-f20ea4146898")
	if err != nil {
		log.Fatalf("failed to parse MessageID: %v", err)
	}
}

func GetPath(file string) string {
	return filepath.Join(rootDir, applicationDir, file)
}

func GetUserPath(phone string, file string) string {
	pp := filepath.Join(rootDir, applicationDir, usersDir, phone, file)
	fmt.Println(pp)
	return pp
}
