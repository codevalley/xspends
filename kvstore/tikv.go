/*
MIT License

# Copyright (c) 2023 Narayan Babu

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
package kvstore

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/tikv/client-go/v2/config"
	"github.com/tikv/client-go/v2/rawkv"
)

const ClientPoolSize = 10
const DefaultMonitoringInterval = 30 * time.Second

var ClientPool chan RawKVClientInterface
var pdAddrs = []string{"tidb-cluster-pd.tidb-cluster.svc.cluster.local:2379"}
var security = config.Security{}

// setupClientPool creates a pool of TiKV clients and returns a channel of clients.
// The size of the pool is determined by the clientPoolSize variable.
// Each client is created using the rawkv.NewClient function with the provided context, PD addresses, and security options.
// If an error occurs while creating a client, the function will return an error.
// The function returns a channel of clients that can be used to perform operations on TiKV.
func setupClientPool(ctx context.Context, useMock bool) (chan RawKVClientInterface, error) {
	clientPool := make(chan RawKVClientInterface, ClientPoolSize)
	for i := 0; i < ClientPoolSize; i++ {
		var client RawKVClientInterface
		if useMock {
			//client = NewMockRawKVClientInterface(nil) // Assuming you have the mock generated
		} else {
			actualClient, err := rawkv.NewClient(ctx, pdAddrs, security)
			if err != nil {
				return nil, fmt.Errorf("failed to create TiKV client: %v", err)
			}
			client = &RawKVClientWrapper{
				client: actualClient,
			}
		}
		clientPool <- client
	}
	return clientPool, nil
}

func GetClientFromPool(clientPool ...chan RawKVClientInterface) RawKVClientInterface {
	var cp chan RawKVClientInterface
	if len(clientPool) > 0 && clientPool[0] != nil {
		cp = clientPool[0]
	} else {
		cp = ClientPool
	}
	if len(cp) > 0 && cap(cp) > 0 {
		return <-cp
	} else {
		return nil
	}

}

func SetupKV(ctx context.Context, useMock bool) {
	var err error
	ClientPool, err = setupClientPool(ctx, useMock) // not mock
	if err != nil {
		log.Fatalf("Failed to create TiKV client: %v", err)
	}
}
