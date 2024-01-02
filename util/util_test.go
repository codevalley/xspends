package util 

import (
    "net/http/httptest"
    "testing"
    "github.com/gin-gonic/gin"
)
func TestContains(t *testing.T) {
    tests := []struct {
        slice []string
        item  string
        want  bool
    }{
        // Test case 1: Item is present in the slice
        {[]string{"apple", "banana", "cherry"}, "banana", true},
        
        // Test case 2: Item is not present in the slice
        {[]string{"apple", "banana", "cherry"}, "mango", false},
        
        // Test case 3: Empty slice
        {[]string{}, "banana", false},
        
        // Add more test cases as needed
    }

    for _, tt := range tests {
        got := Contains(tt.slice, tt.item)
        if got != tt.want {
            t.Errorf("Contains(%v, %v) = %v; want %v", tt.slice, tt.item, got, tt.want)
        }
    }
}

func TestGetUserIDFromQuery(t *testing.T) {
    // Setting up Gin
    gin.SetMode(gin.TestMode)

    tests := []struct {
        queryKey   string
        queryValue string
        wantUserID int64
        wantOk     bool
    }{
        // Test case 1: Valid user ID in query
        {"user", "123", 123, true},

        // Test case 2: Invalid user ID in query (non-integer)
        {"user", "abc", 0, false},

        // Test case 3: No user ID in query
        {"user", "", 0, false},

        // Add more test cases as needed
    }

    for _, tt := range tests {
        // Create a request to pass to our handler.
        req := httptest.NewRequest("GET", "/?"+tt.queryKey+"="+tt.queryValue, nil)

        // We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
        w := httptest.NewRecorder()

        // Create a Gin context
        _, r := gin.CreateTestContext(w)
        r.GET("/", func(c *gin.Context) {
            userID, ok := GetUserIDFromQuery(c, tt.queryKey)

            // Assert
            if userID != tt.wantUserID || ok != tt.wantOk {
                t.Errorf("GetUserIDFromQuery(%v) = (%v, %v); want (%v, %v)", tt.queryKey, userID, ok, tt.wantUserID, tt.wantOk)
            }
        })

        // Serve the HTTP request
        r.ServeHTTP(w, req)
    }
}
