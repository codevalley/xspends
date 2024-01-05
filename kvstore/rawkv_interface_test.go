package kvstore

import (
	"context"
	"testing"

	"xspends/kvstore/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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

func TestRawKVClientInterface_Put(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockRawKVClientInterface(ctrl)

	ctx := context.Background()
	key := []byte("key")
	value := []byte("value")

	mockClient.EXPECT().Put(ctx, key, value).Return(nil).Times(1)

	err := mockClient.Put(ctx, key, value)

	assert.NoError(t, err)
}

func TestRawKVClientInterface_PutWithTTL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockRawKVClientInterface(ctrl)

	ctx := context.Background()
	key := []byte("key")
	value := []byte("value")
	ttl := uint64(100)

	mockClient.EXPECT().PutWithTTL(ctx, key, value, ttl).Return(nil).Times(1)

	err := mockClient.PutWithTTL(ctx, key, value, ttl)

	assert.NoError(t, err)
}

func TestRawKVClientInterface_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockRawKVClientInterface(ctrl)

	ctx := context.Background()
	key := []byte("key")

	mockClient.EXPECT().Delete(ctx, key).Return(nil).Times(1)

	err := mockClient.Delete(ctx, key)

	assert.NoError(t, err)
}

func TestRawKVClientInterface_Scan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockRawKVClientInterface(ctrl)

	ctx := context.Background()
	startKey := []byte("startKey")
	endKey := []byte("endKey")
	limit := 10
	keys := [][]byte{[]byte("key1"), []byte("key2")}
	values := [][]byte{[]byte("value1"), []byte("value2")}

	mockClient.EXPECT().Scan(ctx, startKey, endKey, limit).Return(keys, values, nil).Times(1)

	resultKeys, resultValues, err := mockClient.Scan(ctx, startKey, endKey, limit)

	assert.NoError(t, err)
	assert.Equal(t, keys, resultKeys)
	assert.Equal(t, values, resultValues)
}

func TestRawKVClientWrapper_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockRawKVClientInterface(ctrl)
	wrapper := NewRawKVClientWrapper(mockClient)

	ctx := context.Background()
	key := []byte("key")
	value := []byte("value")

	// Setting up expectations
	mockClient.EXPECT().Get(ctx, key).Return(value, nil).Times(1)

	// Call the method we're testing
	result, err := wrapper.Get(ctx, key)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, value, result)
}

func TestRawKVClientWrapper_Put(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockRawKVClientInterface(ctrl)
	wrapper := NewRawKVClientWrapper(mockClient)

	ctx := context.Background()
	key := []byte("key")
	value := []byte("value")

	// Expectation: mockClient.Put is called once and returns nil
	mockClient.EXPECT().Put(ctx, key, value).Return(nil).Times(1)

	// Act
	err := wrapper.Put(ctx, key, value)

	// Assert
	assert.NoError(t, err)
}

func TestRawKVClientWrapper_PutWithTTL(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockRawKVClientInterface(ctrl)
	wrapper := NewRawKVClientWrapper(mockClient)

	ctx := context.Background()
	key := []byte("key")
	value := []byte("value")
	ttl := uint64(100)

	// Expectation: mockClient.PutWithTTL is called once and returns nil
	mockClient.EXPECT().PutWithTTL(ctx, key, value, ttl).Return(nil).Times(1)

	// Act
	err := wrapper.PutWithTTL(ctx, key, value, ttl)

	// Assert
	assert.NoError(t, err)
}

func TestRawKVClientWrapper_Delete(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockRawKVClientInterface(ctrl)
	wrapper := NewRawKVClientWrapper(mockClient)

	ctx := context.Background()
	key := []byte("key")

	// Expectation: mockClient.Delete is called once and returns nil
	mockClient.EXPECT().Delete(ctx, key).Return(nil).Times(1)

	// Act
	err := wrapper.Delete(ctx, key)

	// Assert
	assert.NoError(t, err)
}

func TestRawKVClientWrapper_Scan(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockRawKVClientInterface(ctrl)
	wrapper := NewRawKVClientWrapper(mockClient)

	ctx := context.Background()
	startKey := []byte("startKey")
	endKey := []byte("endKey")
	limit := 10
	keys := [][]byte{[]byte("key1"), []byte("key2")}
	values := [][]byte{[]byte("value1"), []byte("value2")}

	// Expectation: mockClient.Scan is called once and returns keys and values
	mockClient.EXPECT().Scan(ctx, startKey, endKey, limit).Return(keys, values, nil).Times(1)

	// Act
	resultKeys, resultValues, err := wrapper.Scan(ctx, startKey, endKey, limit)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, keys, resultKeys)
	assert.Equal(t, values, resultValues)
}
