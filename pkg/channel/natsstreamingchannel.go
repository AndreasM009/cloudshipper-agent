package channel

import (
	snatsio "github.com/nats-io/stan.go"
)

// NatsStreamingChannel streaming channel
type NatsStreamingChannel struct {
	SnatConnection       *NatsStreamingConnection
	SnatNativeConnection snatsio.Conn
	NatsPublishName      string
	isPooledConnection   bool
}

// NewNatsStreamingChannel new instance
func NewNatsStreamingChannel(channelName string, natsServerURLs []string, natsConnectionName string, clusterID string, clientID string) (*NatsStreamingChannel, error) {
	con, err := NewNatsStreamingConnection(natsServerURLs, natsConnectionName, clusterID, clientID)
	if err != nil {
		return nil, err
	}

	return &NatsStreamingChannel{
		NatsPublishName:      channelName,
		SnatConnection:       con,
		SnatNativeConnection: con.SnatConnection,
		isPooledConnection:   false,
	}, nil
}

// NewNatsStreamingChannelFromPool new channel from connection pool
func NewNatsStreamingChannelFromPool(channelName, clusterID, clientID string) (*NatsStreamingChannel, error) {
	pool := GetNatsStreamingConnectionPoolInstance()
	c, err := pool.Get(clusterID, clientID)

	if err != nil {
		return nil, err
	}

	return &NatsStreamingChannel{
		NatsPublishName:      channelName,
		SnatConnection:       c,
		SnatNativeConnection: c.SnatConnection,
		isPooledConnection:   true,
	}, nil
}

// Close closes connections
func (channel *NatsStreamingChannel) Close() {
	if !channel.isPooledConnection {
		channel.SnatConnection.Close()
	}
}
