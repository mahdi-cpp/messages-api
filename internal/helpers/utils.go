package helpers

import "github.com/gin-gonic/gin"

// GetUserID از Gin context، user_id را به صورت string دریافت می‌کند.
func GetUserID(c *gin.Context) (string, bool) {

	//// این تابع باید بعد از middleware احراز هویت استفاده شود
	//userID, exists := c.Get("X-User-ID")
	//if !exists {
	//	return "", false
	//}
	//
	//userIDStr, ok := userID.(string)
	//if !ok {
	//	return "", false
	//}
	//
	//return userIDStr, true

	//// این تابع باید بعد از middleware احراز هویت استفاده شود
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		return "", false
	}

	return userID, true
}
