package channel

import snatsio "github.com/nats-io/stan.go"

// NatsStreamingConnection connection to Nats Streaming Server
type NatsStreamingConnection struct {
	natsConnection *NatsConnection
	clientID       string
	clusterID      string
	SnatConnection snatsio.Conn
}

// NewNatsStreamingConnection new nats streaming server connection
func NewNatsStreamingConnection(natsServerURLs []string, natsConnectionName string, clusterID string, clientID string) (*NatsStreamingConnection, error) {
	nc, err := NewNatsConnection(natsServerURLs, natsConnectionName)
	if err != nil {
		return nil, err
	}

	sc, err := snatsio.Connect(clusterID, clientID, snatsio.NatsConn(nc.NatsConn))
	if err != nil {
		return nil, err
	}

	return &NatsStreamingConnection{
		natsConnection: nc,
		clientID:       clientID,
		clusterID:      clusterID,
		SnatConnection: sc,
	}, nil
}

// NewNatsStreamingConnectionWithPooledConnection with nats connection from pool
func NewNatsStreamingConnectionWithPooledConnection(connectionName, clusterID, clientID string) (*NatsStreamingConnection, error) {
	c, err := GetNatsConnectionPoolInstance().Get(connectionName)
	if err != nil {
		return nil, err
	}

	sc, err := snatsio.Connect(clusterID, clientID, snatsio.NatsConn(c.NatsConn))
	if err != nil {
		return nil, err
	}

	return &NatsStreamingConnection{
		natsConnection: c,
		clientID:       clientID,
		clusterID:      clusterID,
		SnatConnection: sc,
	}, nil
}

// Close closes connection
func (connection *NatsStreamingConnection) Close() {
	connection.SnatConnection.Close()
	connection.natsConnection.Close()
}
