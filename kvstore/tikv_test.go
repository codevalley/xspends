package kvstore

import (
	"context"
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

	// Code under test
	clientPool, err := setupClientPool(ctx, useMock)
	if err != nil {
		t.Fatalf("Failed to create client pool: %v", err)
	}

	// Assertions
	for i := 0; i < ClientPoolSize; i++ {
		client, more := <-clientPool
		if !more {
			t.Fatalf("Expected more clients in the pool, but channel is closed")
		}

		_, ok := client.(*mock.MockRawKVClientInterface) //not sure if this is correct
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

func TestCreateMockClientIfUseMockIsTrue(t *testing.T) {
	// Setup
	ctx := context.Background()
	useMock := true
	const expectedClientCount = 1 // Adjust this to match the expected number of clients

	// Mock setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := mock.NewMockRawKVClientInterface(ctrl)

	// Expectations
	// Adjust the expectation to match how many times you expect the method to be called
	mockClient.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil, nil).AnyTimes()

	// Code under test
	clientPool, err := setupClientPool(ctx, useMock)
	if err != nil {
		t.Fatalf("Failed to create client pool: %v", err)
	}

	// Assertions
	count := 0
	for client := range clientPool {
		_, isMockClient := client.(*mock.MockRawKVClientInterface)
		if !isMockClient {
			t.Errorf("Expected client to be of type *mock.MockRawKVClientInterface, but got %T", client)
		}

		count++
		if count >= expectedClientCount {
			break // Exit the loop after processing the expected number of clients
		}
	}

	if count != expectedClientCount {
		t.Errorf("Expected %d clients in the pool, but got %d", expectedClientCount, count)
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

// GetClientFromPool returns a client from the pool if there are any available
func TestGetClientFromPoolReturnsClientIfAvailable(t *testing.T) {
	// Setup
	ctx := context.Background()
	useMock := true

	// Mock setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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
