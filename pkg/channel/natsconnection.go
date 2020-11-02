package channel

import (
	"strings"
	"time"

	natsio "github.com/nats-io/nats.go"
)

// NatsConnection connection to nats server
type NatsConnection struct {
	natsURLs           []string
	NatsConn           *natsio.Conn
	natsConnectionName string
}

// NewNatsConnection creates new connection
func NewNatsConnection(natsServerURLs []string, natsConnectionName string) (*NatsConnection, error) {
	nc, err := natsio.Connect(
		strings.Join(natsServerURLs, ","), natsio.Name(natsConnectionName), natsio.Timeout(10*time.Second),
		natsio.PingInterval(20*time.Second), natsio.MaxPingsOutstanding(5), natsio.ReconnectBufSize(10*1024*1024), natsio.NoEcho())
	if err != nil {
		return nil, err
	}

	return &NatsConnection{
		natsURLs:           natsServerURLs,
		NatsConn:           nc,
		natsConnectionName: natsConnectionName,
	}, nil
}

// Close closes connection
func (connection *NatsConnection) Close() {
	connection.NatsConn.Close()
}
