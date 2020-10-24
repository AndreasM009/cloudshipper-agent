package channel

import (
	"strings"
	"time"

	natsio "github.com/nats-io/nats.go"
	snatsio "github.com/nats-io/stan.go"
)

// NatsStreamingChannel streaming channel
type NatsStreamingChannel struct {
	natsConnection     *natsio.Conn
	SnatConnection     snatsio.Conn
	NatsPublishName    string
	natsURLs           []string
	natsConnectionName string
	clientID           string
	clusterID          string
}

// NewNatsStreamingChannel new instance
func NewNatsStreamingChannel(channelName string, natsServerURLs []string, natsConnectionName string, clusterID string, clientID string) (*NatsStreamingChannel, error) {
	nc, err := natsio.Connect(
		strings.Join(natsServerURLs, ","), natsio.Name(natsConnectionName), natsio.Timeout(10*time.Second),
		natsio.PingInterval(20*time.Second), natsio.MaxPingsOutstanding(5), natsio.ReconnectBufSize(10*1024*1024), natsio.NoEcho())
	if err != nil {
		return nil, err
	}

	sc, err := snatsio.Connect(clusterID, clientID, snatsio.NatsConn(nc))
	if err != nil {
		return nil, err
	}

	return &NatsStreamingChannel{
		NatsPublishName:    channelName,
		natsURLs:           natsServerURLs,
		natsConnectionName: natsConnectionName,
		SnatConnection:     sc,
		natsConnection:     nc,
		clientID:           clientID,
		clusterID:          clusterID,
	}, nil
}

// Close closes connections
func (channel *NatsStreamingChannel) Close() {
	channel.SnatConnection.Close()
	channel.natsConnection.Close()
}
