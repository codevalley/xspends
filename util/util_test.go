import (
    "testing"
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
