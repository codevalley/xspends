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
