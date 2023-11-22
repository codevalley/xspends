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
	"errors"
	"log"
	"sync"

	"github.com/sony/sonyflake"
)

var (
	snowflakeGenerator *sonyflake.Sonyflake
	once               sync.Once
)

// InitializeSnowflake initializes the Snowflake ID generator.
// Call this function once when your application starts.
func InitializeSnowflake() {
	once.Do(func() {
		snowflakeGenerator = sonyflake.NewSonyflake(sonyflake.Settings{})
		if snowflakeGenerator == nil {
			log.Fatal("Failed to initialize Snowflake generator")
		}
	})
}

// GenerateSnowflakeID generates a new Snowflake ID.
func GenerateSnowflakeID() (int64, error) {
	if snowflakeGenerator == nil {
		return 0, errors.New("Snowflake generator is not initialized")
	}

	id, err := snowflakeGenerator.NextID()
	if err != nil {
		return 0, err
	}
	return int64(id), nil
}
