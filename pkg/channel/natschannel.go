package channel

import (
	natsio "github.com/nats-io/nats.go"
)

// NatsChannel message channel based on nats.io
type NatsChannel struct {
	NatsPublishName  string
	NatsNativeConn   *natsio.Conn
	NatsConnection   *NatsConnection
	isPoolConnection bool
}

// NewNatsChannel new instance
func NewNatsChannel(channelName string, natsServerURLs []string, natsConnectionName string) (*NatsChannel, error) {
	con, err := NewNatsConnection(natsServerURLs, natsConnectionName)
	if err != nil {
		return nil, err
	}

	return &NatsChannel{
		NatsPublishName:  channelName,
		NatsConnection:   con,
		NatsNativeConn:   con.NatsConn,
		isPoolConnection: false,
	}, nil
}

// NewNatsChannelFromPool new channel with connection from pool
func NewNatsChannelFromPool(channelName, connectionName string) (*NatsChannel, error) {
	pool := GetNatsConnectionPoolInstance()
	c, err := pool.Get(connectionName)

	if err != nil {
		return nil, err
	}

	return &NatsChannel{
		NatsPublishName:  channelName,
		NatsConnection:   c,
		NatsNativeConn:   c.NatsConn,
		isPoolConnection: true,
	}, nil
}

// Close closes the channel
func (channel *NatsChannel) Close() {
	if !channel.isPoolConnection {
		channel.NatsConnection.Close()
	}
}
