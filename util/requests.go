/*
MIT License

Copyright (c) 2023 Narayan Babu

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

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
