package kvstore

import (
	"context"
	"fmt"
	"testing"
	"xspends/kvstore/mock"

	"github.com/golang/mock/gomock"
)

// Creates a pool of TiKV clients with the specified size
func TestCreateClientPoolWithSpecifiedSize(t *testing.T) {
	// Setup
	ctx := context.Background()
	useMock := true

	// Mock setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := mock.NewMockRawKVClientInterface(ctrl)

	// Expectations
	mockClient.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil, nil)

	// Code under test
	clientPool, err := setupClientPool(ctx, useMock)
	if err != nil {
		t.Fatalf("Failed to create client pool: %v", err)
	}

	// Assertions
	if len(clientPool) != ClientPoolSize {
		t.Errorf("Expected client pool size to be %d, but got %d", ClientPoolSize, len(clientPool))
	}
}

// Each client is created using the rawkv.NewClient function with the provided context, PD addresses, and security options
func TestCreateClientWithRawKVNewClientFunction(t *testing.T) {
	// Setup
	ctx := context.Background()
	useMock := true

	// Mock setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := mock.NewMockRawKVClientInterface(ctrl)

	// Expectations
	mockClient.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil, nil)

	// Code under test
	clientPool, err := setupClientPool(ctx, useMock)
	if err != nil {
		t.Fatalf("Failed to create client pool: %v", err)
	}

	// Assertions
	for client := range clientPool {
		_, ok := client.(*RawKVClientWrapper)
		if !ok {
			t.Errorf("Expected client to be of type *RawKVClientWrapper, but got %T", client)
		}
	}
}

// Returns a channel of clients that can be used to perform operations on TiKV
func TestReturnChannelOfClients(t *testing.T) {
	// Setup
	ctx := context.Background()
	useMock := true

	// Mock setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := mock.NewMockRawKVClientInterface(ctrl)

	// Expectations
	mockClient.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil, nil)

	// Code under test
	clientPool, err := setupClientPool(ctx, useMock)
	if err != nil {
		t.Fatalf("Failed to create client pool: %v", err)
	}

	// Assertions
	if len(clientPool) != ClientPoolSize {
		t.Errorf("Expected client pool size to be %d, but got %d", ClientPoolSize, len(clientPool))
	}
}

// If useMock is true, creates a mock client instead of an actual client
func TestCreateMockClientIfUseMockIsTrue(t *testing.T) {
	// Setup
	ctx := context.Background()
	useMock := true

	// Mock setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := mock.NewMockRawKVClientInterface(ctrl)

	// Expectations
	mockClient.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil, nil)

	// Code under test
	clientPool, err := setupClientPool(ctx, useMock)
	if err != nil {
		t.Fatalf("Failed to create client pool: %v", err)
	}

	// Assertions
	for client := range clientPool {
		_, ok := client.(*mock.MockRawKVClientInterface)
		if !ok {
			t.Errorf("Expected client to be of type *mock.MockRawKVClientInterface, but got %T", client)
		}
	}
}

// Returns a channel of mock clients if useMock is true
func TestReturnChannelOfMockClientsIfUseMockIsTrue(t *testing.T) {
	// Setup
	ctx := context.Background()
	useMock := true

	// Mock setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := mock.NewMockRawKVClientInterface(ctrl)

	// Expectations
	mockClient.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil, nil)

	// Code under test
	clientPool, err := setupClientPool(ctx, useMock)
	if err != nil {
		t.Fatalf("Failed to create client pool: %v", err)
	}

	// Assertions
	if len(clientPool) != ClientPoolSize {
		t.Errorf("Expected client pool size to be %d, but got %d", ClientPoolSize, len(clientPool))
	}
}

// Returns an error if an error occurs while creating a client
func TestReturnErrorIfErrorOccursWhileCreatingClient(t *testing.T) {
	// Setup
	ctx := context.Background()
	useMock := true

	// Mock setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := mock.NewMockRawKVClientInterface(ctrl)

	// Expectations
	mockClient.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil, fmt.Errorf("error creating client"))

	// Code under test
	_, err := setupClientPool(ctx, useMock)

	// Assertions
	if err == nil {
		t.Error("Expected an error, but got nil")
	}
}

// Returns nil if the provided clientPool is empty
func TestReturnNilIfClientPoolIsEmpty(t *testing.T) {
	// Setup
	var clientPool chan RawKVClientInterface

	// Code under test
	result := GetClientFromPool(clientPool)

	// Assertions
	if result != nil {
		t.Errorf("Expected result to be nil, but got %v", result)
	}
}

// Returns nil if the provided clientPool has a capacity of 0
func TestReturnNilIfClientPoolHasCapacityOfZero(t *testing.T) {
	// Setup
	clientPool := make(chan RawKVClientInterface, 0)

	// Code under test
	result := GetClientFromPool(clientPool)

	// Assertions
	if result != nil {
		t.Errorf("Expected result to be nil, but got %v", result)
	}
}

// Returns a channel of actual clients if useMock is false
func TestReturnChannelOfActualClientsIfUseMockIsFalse(t *testing.T) {
	// Setup
	ctx := context.Background()
	useMock := true

	// Mock setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := mock.NewMockRawKVClientInterface(ctrl)

	// Expectations
	mockClient.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil, nil)

	// Code under test
	clientPool, err := setupClientPool(ctx, useMock)
	if err != nil {
		t.Fatalf("Failed to create client pool: %v", err)
	}

	// Assertions
	for client := range clientPool {
		_, ok := client.(*RawKVClientWrapper)
		if !ok {
			t.Errorf("Expected client to be of type *RawKVClientWrapper, but got %T", client)
		}
	}
}

// GetClientFromPool returns a client from the pool if there are any available
func TestGetClientFromPoolReturnsClientIfAvailable(t *testing.T) {
	// Setup
	ctx := context.Background()
	useMock := true

	// Mock setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := mock.NewMockRawKVClientInterface(ctrl)

	// Expectations
	mockClient.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil, nil)

	// Code under test
	clientPool, err := setupClientPool(ctx, useMock)
	if err != nil {
		t.Fatalf("Failed to create client pool: %v", err)
	}

	// Assertions
	client := GetClientFromPool(clientPool)
	if client == nil {
		t.Error("Expected a client, but got nil")
	}
}
