package channel

import (
	"errors"
	"fmt"
	"sync"
)

var theNatsConnectionPool *NatsConnectionPool
var onceNatsConnectionPool *sync.Once = &sync.Once{}
var natsConnectionPoolMutex *sync.Mutex = &sync.Mutex{}

// NatsConnectionPool the pool
type NatsConnectionPool struct {
	pool map[string]*NatsConnection
}

// GetNatsConnectionPoolInstance gets the instance
func GetNatsConnectionPoolInstance() *NatsConnectionPool {
	onceNatsConnectionPool.Do(func() {
		theNatsConnectionPool = &NatsConnectionPool{
			pool: make(map[string]*NatsConnection),
		}
	})

	return theNatsConnectionPool
}

// Add adds connection to the pool
func (pool *NatsConnectionPool) Add(connectionName string, connection *NatsConnection) error {
	if connectionName == "" {
		return errors.New("connectionName can not be empty")
	}

	if connection == nil {
		return errors.New("connection can not be nil")
	}

	natsConnectionPoolMutex.Lock()
	defer natsConnectionPoolMutex.Unlock()

	if _, ok := pool.pool[connectionName]; ok {
		return errors.New("a connection with given name already in the pool")
	}

	pool.pool[connectionName] = connection
	return nil
}

// Get gets connection from pool
func (pool *NatsConnectionPool) Get(connectionName string) (*NatsConnection, error) {
	if connectionName == "" {
		return nil, errors.New("connectionName can not be empty")
	}

	natsConnectionPoolMutex.Lock()
	defer natsConnectionPoolMutex.Unlock()

	if c, ok := pool.pool[connectionName]; ok {
		return c, nil
	}

	return nil, fmt.Errorf("connection %s not in pool", connectionName)
}
