package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	ErrorUserID         = "User ID is not a string"
	ErrorInvalidRequest = "Invalid request"
)

// AbortWithError یک پاسخ JSON خطا را ارسال و درخواست را Abort می‌کند.
func AbortWithError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
	c.Abort()
}

// AbortWithUserIDInvalid یک پاسخ خطا برای زمانی که user_id نامعتبر است، ارسال می‌کند.
func AbortWithUserIDInvalid(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": ErrorUserID})
	c.Abort() // This stops the next handler from running
}

func AbortWithRequestInvalid(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": ErrorInvalidRequest})
	c.Abort()
}
