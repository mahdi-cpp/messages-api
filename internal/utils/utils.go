package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
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
