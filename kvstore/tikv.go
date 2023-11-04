package kvstore

import (
	"context"
	"log"
	"time"

	"github.com/tikv/client-go/v2/config"
	"github.com/tikv/client-go/v2/rawkv"
)

const ClientPoolSize = 10
const DefaultMonitoringInterval = 30 * time.Second

var ClientPool chan RawKVClientInterface
var ctx = context.Background()
var pdAddrs = []string{"tidb-cluster-pd.tidb-cluster.svc.cluster.local:2379"}
var security = config.Security{}

// setupClientPool creates a pool of TiKV clients and returns a channel of clients.
// The size of the pool is determined by the clientPoolSize variable.
// Each client is created using the rawkv.NewClient function with the provided context, PD addresses, and security options.
// If an error occurs while creating a client, the function will log a fatal error and exit.
// The function returns a channel of clients that can be used to perform operations on TiKV.
func setupClientPool(useMock bool) chan RawKVClientInterface {
	clientPool := make(chan RawKVClientInterface, ClientPoolSize)
	for i := 0; i < ClientPoolSize; i++ {
		var client RawKVClientInterface
		if useMock {
			//client = NewMockRawKVClientInterface(nil) // Assuming you have the mock generated
		} else {
			actualClient, err := rawkv.NewClient(ctx, pdAddrs, security)
			if err != nil {
				log.Fatalf("Failed to create TiKV client: %v", err)
			}
			client = &RawKVClientWrapper{
				client: actualClient,
			}
		}
		clientPool <- client
	}
	return clientPool
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

func SetupKV(useMock bool) {
	// setupLogging(LogFile)
	ClientPool = setupClientPool(false) // not mock
	// setupMonitoring(clientPool)
}
