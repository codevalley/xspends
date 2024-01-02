package kvstore

import (
    "context"
    "testing"

    "github.com/golang/mock/gomock"
    "github.com/stretchr/testify/assert"
    "xspends/kvstore/mock" 
)

func TestRawKVClientInterface_Get(t *testing.T) {
    // Initialize the mock controller
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    // Create an instance of the mock RawKVClientInterface
    mockClient := mock.NewMockRawKVClientInterface(ctrl)

    // Setting up expectations
    ctx := context.Background()
    key := []byte("key")
    value := []byte("value")
    mockClient.EXPECT().Get(ctx, key).Return(value, nil).Times(1)

    // Call the method we're testing
    result, err := mockClient.Get(ctx, key)

    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, value, result)
}

// Similar test functions for Put, PutWithTTL, Delete, Scan...

