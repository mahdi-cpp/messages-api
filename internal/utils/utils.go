package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"os"
	"strconv"
)

// utils-------------------------------------

func GetFileSize(filepath string) (int64, error) {
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

func GetUserId(c *gin.Context) (int, error) {

	userIDStr := c.GetHeader("userID")
	fmt.Println(userIDStr)

	return strconv.Atoi(userIDStr)
}

func GenerateUUID() (string, error) {
	u7, err2 := uuid.NewV7()
	if err2 != nil {
		return "", fmt.Errorf("error generating UUIDv7: %w", err2)
	}
	return u7.String(), nil
}
