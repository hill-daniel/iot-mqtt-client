package mqtt

import (
	"crypto/tls"
	mq "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
	pb "github.com/hill-daniel/iot-mqtt-client/proto"
	"github.com/pkg/errors"
)

// Publisher publishes device data via MQTT
type Publisher struct {
	client mq.Client
}

// New creates a new MQTT Publisher
func New(brokerURL string, clientID string, tlsConfig *tls.Config) *Publisher {
	opts := mq.NewClientOptions().AddBroker(brokerURL)
	opts.SetClientID(clientID)
	opts.SetTLSConfig(tlsConfig)
	client := mq.NewClient(opts)
	return &Publisher{client: client}
}

// Publish publishes the given DeviceMessage to the topic
func (p *Publisher) Publish(topic string, device *pb.Device) error {
	if !p.client.IsConnected() {
		token := p.client.Connect()

		if token.Wait() && token.Error() != nil {
			return token.Error()
		}
	}

	protoDevice, err := proto.Marshal(device)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal device %v", device)
	}

	// Quality of Service (QoS) of 0 means message is delivered at most once, which is ok for our case right now
	token := p.client.Publish(topic, 0, false, protoDevice)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
