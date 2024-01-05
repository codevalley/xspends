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

func TestKVGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKVClient := mock.NewMockRawKVClientInterface(ctrl)

	// Set up expected calls on the mock
	testKey := []byte("key")
	testValue := []byte("value")
	mockKVClient.EXPECT().Get(gomock.Any(), testKey).Return(testValue, nil).Times(1)

	// Call the method under test
	value, err := mockKVClient.Get(context.Background(), testKey)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, testValue, value)
}

func TestKVPut(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKVClient := mock.NewMockRawKVClientInterface(ctrl)

	// Set up expected calls on the mock
	testKey := []byte("key")
	testValue := []byte("value")
	mockKVClient.EXPECT().Put(gomock.Any(), testKey, testValue).Return(nil).Times(1)

	// Call the method under test
	err := mockKVClient.Put(context.Background(), testKey, testValue)

	// Assertions
	assert.NoError(t, err)
}

func TestKVPutWithTTL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKVClient := mock.NewMockRawKVClientInterface(ctrl)

	// Set up expected calls on the mock
	testKey := []byte("key")
	testValue := []byte("value")
	testTTL := uint64(100)
	mockKVClient.EXPECT().PutWithTTL(gomock.Any(), testKey, testValue, testTTL).Return(nil).Times(1)

	// Call the method under test
	err := mockKVClient.PutWithTTL(context.Background(), testKey, testValue, testTTL)

	// Assertions
	assert.NoError(t, err)
}

func TestKVDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKVClient := mock.NewMockRawKVClientInterface(ctrl)

	// Set up expected calls on the mock
	testKey := []byte("key")
	mockKVClient.EXPECT().Delete(gomock.Any(), testKey).Return(nil).Times(1)

	// Call the method under test
	err := mockKVClient.Delete(context.Background(), testKey)

	// Assertions
	assert.NoError(t, err)
}

func TestKVScan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKVClient := mock.NewMockRawKVClientInterface(ctrl)

	// Set up expected calls on the mock
	startKey := []byte("start")
	endKey := []byte("end")
	limit := 10
	expectedKeys := [][]byte{[]byte("key1"), []byte("key2")}
	expectedValues := [][]byte{[]byte("value1"), []byte("value2")}
	mockKVClient.EXPECT().Scan(gomock.Any(), startKey, endKey, limit).Return(expectedKeys, expectedValues, nil).Times(1)

	// Call the method under test
	keys, values, err := mockKVClient.Scan(context.Background(), startKey, endKey, limit)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedKeys, keys)
	assert.Equal(t, expectedValues, values)
}
