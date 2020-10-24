package channel

import (
	"strings"
	"time"

	natsio "github.com/nats-io/nats.go"
)

// NatsChannel message channel based on nats.io
type NatsChannel struct {
	NatsPublishName    string
	natsURLs           []string
	NatsConn           *natsio.Conn
	natsConnectionName string
}

// NewNatsChannel new instance
func NewNatsChannel(channelName string, natsServerURLs []string, natsConnectionName string) (*NatsChannel, error) {
	nc, err := natsio.Connect(
		strings.Join(natsServerURLs, ","), natsio.Name(natsConnectionName), natsio.Timeout(10*time.Second),
		natsio.PingInterval(20*time.Second), natsio.MaxPingsOutstanding(5), natsio.ReconnectBufSize(10*1024*1024), natsio.NoEcho())
	if err != nil {
		return nil, err
	}

	return &NatsChannel{
		NatsPublishName:    channelName,
		natsURLs:           natsServerURLs,
		natsConnectionName: natsConnectionName,
		NatsConn:           nc,
	}, nil
}

// Close closes the channel
func (channel *NatsChannel) Close() {
	channel.NatsConn.Close()
}
