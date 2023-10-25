package util

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetFloatFromQuery retrieves a float value from query parameters with a default value
func GetFloatFromQuery(c *gin.Context, key string, defaultValue float64) float64 {
	value, err := strconv.ParseFloat(c.DefaultQuery(key, strconv.FormatFloat(defaultValue, 'f', -1, 64)), 64)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetIntFromQuery retrieves an integer value from query parameters with a default value
func GetIntFromQuery(c *gin.Context, key string, defaultValue int) int {
	value, err := strconv.Atoi(c.DefaultQuery(key, strconv.Itoa(defaultValue)))
	if err != nil {
		return defaultValue
	}
	return value
}

func GetUserIDFromQuery(c *gin.Context, key string) (int64, bool) {
	userID, err := strconv.ParseInt(c.Query(key), 10, 64)
	if err != nil {
		return 0, false
	}
	return userID, true
}
