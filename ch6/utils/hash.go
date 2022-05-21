package utils

import uuid "github.com/satori/go.uuid"

func GenUUID() string {
	u1 := uuid.NewV4().String()

	return u1
}

func CheckSign(appId string, appSecret string, token string) bool {

	return true
}
