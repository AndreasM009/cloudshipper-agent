package channel

import (
	"errors"
	"fmt"
	"sync"
)

var theNatsStreamingConnectionPool *NatsStreamingConnectionPool
var onceNatsStreamingConnectionPool *sync.Once = &sync.Once{}
var natsStreamingConnectionPoolMutex *sync.Mutex = &sync.Mutex{}

// NatsStreamingConnectionPool pool for nats streaming server connections
type NatsStreamingConnectionPool struct {
	pool map[string]*NatsStreamingConnection
}

// GetNatsStreamingConnectionPoolInstance gets the pool instance
func GetNatsStreamingConnectionPoolInstance() *NatsStreamingConnectionPool {
	onceNatsStreamingConnectionPool.Do(func() {
		theNatsStreamingConnectionPool = &NatsStreamingConnectionPool{
			pool: make(map[string]*NatsStreamingConnection),
		}
	})

	return theNatsStreamingConnectionPool
}

// Add adds connection to the pool
func (pool *NatsStreamingConnectionPool) Add(clusterID, clientID string, connection *NatsStreamingConnection) error {
	if clusterID == "" {
		return errors.New("clusterID can not be empty")
	}

	if clientID == "" {
		return errors.New("clientID ac not be empty")
	}

	if connection == nil {
		return errors.New("connection can not be nil")
	}

	key := fmt.Sprintf("%s-%s", clusterID, clientID)

	natsStreamingConnectionPoolMutex.Lock()
	defer natsStreamingConnectionPoolMutex.Unlock()

	if _, ok := pool.pool[key]; ok {
		return errors.New("a connection with the given clientID and clusterID already exists")
	}

	pool.pool[key] = connection
	return nil
}

// Get gets connection from pool
func (pool *NatsStreamingConnectionPool) Get(clusterID, clientID string) (*NatsStreamingConnection, error) {
	if clusterID == "" {
		return nil, errors.New("clusterID can not be empty")
	}

	if clientID == "" {
		return nil, errors.New("clientID ac not be empty")
	}

	key := fmt.Sprintf("%s-%s", clusterID, clientID)

	natsStreamingConnectionPoolMutex.Lock()
	defer natsStreamingConnectionPoolMutex.Unlock()

	if c, ok := pool.pool[key]; ok {
		return c, nil
	}

	return nil, fmt.Errorf("connection with key %s not found in pool", key)
}
