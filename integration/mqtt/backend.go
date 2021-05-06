package mqtt

import (
	"fmt"
	"sync"
	"text/template"
	"time"

	"github.com/apex/log"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/h3c/mygolang/config"
	"github.com/h3c/mygolang/integration/mqtt/auth"
	"github.com/pkg/errors"
)

// Backend implements a MQTT backend.
type Backend struct {
	auth auth.Authentication

	conn       paho.Client
	connMux    sync.RWMutex
	connClosed bool
	clientOpts *paho.ClientOptions

	terminateOnConnectError bool

	qos                  uint8
	eventTopicTemplate   *template.Template
	commandTopicTemplate *template.Template
}

// NewBackend creates a new Backend.
func NewBackend(conf config.Config) (*Backend, error) {
	var err error

	b := Backend{
		qos:                     conf.Integration.MQTT.Auth.Generic.QOS,
		terminateOnConnectError: conf.Integration.MQTT.TerminateOnConnectError,
		clientOpts:              paho.NewClientOptions(),
	}

	fmt.Println("NewBackend")
	switch conf.Integration.MQTT.Auth.Type {
	case "generic":
		fmt.Println("generic")
		b.auth, err = auth.NewGenericAuthentication(conf)
		if err != nil {
			return nil, errors.Wrap(err, "integation/mqtt: new generic authentication error")
		}
	default:
		return nil, fmt.Errorf("integration/mqtt: unknown auth type: %s", conf.Integration.MQTT.Auth.Type)
	}

	b.eventTopicTemplate, err = template.New("event").Parse(conf.Integration.MQTT.EventTopicTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "integration/mqtt: parse event-topic template error")
	}

	b.commandTopicTemplate, err = template.New("event").Parse(conf.Integration.MQTT.CommandTopicTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "integration/mqtt: parse event-topic template error")
	}

	b.clientOpts.SetProtocolVersion(4)
	b.clientOpts.SetAutoReconnect(true) // this is required for buffering messages in case offline!
	b.clientOpts.SetOnConnectHandler(b.onConnected)
	b.clientOpts.SetConnectionLostHandler(b.onConnectionLost)
	b.clientOpts.SetKeepAlive(conf.Integration.MQTT.KeepAlive)
	b.clientOpts.SetMaxReconnectInterval(conf.Integration.MQTT.MaxReconnectInterval)

	if err = b.auth.Init(b.clientOpts); err != nil {
		return nil, errors.Wrap(err, "mqtt: init authentication error")
	}

	return &b, nil
}

// Start starts the integration.
func (b *Backend) Start() error {
	fmt.Println("Start")
	b.connectLoop()
	go b.reconnectLoop()
	// go b.subscribeLoop()
	return nil
}

// Stop stops the integration.
func (b *Backend) Stop() error {
	b.connMux.Lock()
	defer b.connMux.Unlock()

	b.conn.Disconnect(250)
	b.connClosed = true
	return nil
}

func (b *Backend) connect() error {
	fmt.Println("connect")
	b.connMux.Lock()
	defer b.connMux.Unlock()

	if err := b.auth.Update(b.clientOpts); err != nil {
		return errors.Wrap(err, "integration/mqtt: update authentication error")
	}

	b.conn = paho.NewClient(b.clientOpts)
	if token := b.conn.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("connect mqtt failed")
		return token.Error()
	}
	fmt.Println("connect mqtt success: ")

	return nil
}

// connectLoop blocks until the client is connected
func (b *Backend) connectLoop() {
	for {
		if err := b.connect(); err != nil {
			log.WithError(err).Error("integration/mqtt: connection error")
			time.Sleep(time.Second * 2)
		} else {
			break
		}
	}
}

func (b *Backend) disconnect() error {
	b.connMux.Lock()
	defer b.connMux.Unlock()

	b.conn.Disconnect(250)
	return nil
}

func (b *Backend) reconnectLoop() {
	if b.auth.ReconnectAfter() > 0 {
		for {
			b.connMux.RLock()
			closed := b.connClosed
			b.connMux.RUnlock()

			if closed {
				break
			}
			time.Sleep(b.auth.ReconnectAfter())
			log.Info("mqtt: re-connect triggered")

			b.disconnect()
			b.connectLoop()
		}
	}
}

func (b *Backend) onConnected(c paho.Client) {
	log.Info("integration/mqtt: connected to mqtt broker")
}

func (b *Backend) onConnectionLost(c paho.Client, err error) {
	fmt.Println("mqtt lost connect")
	if b.terminateOnConnectError {
		log.Fatal(err.Error())
	}
	log.WithError(err).Error("mqtt: connection error")
}

// func (b *Backend) publish(msg proto.Message) error {
// 	topic := bytes.NewBuffer(nil)

// 	bytes, err := json.Marshal(msg)
// 	if err != nil {
// 		return errors.Wrap(err, "marshal message error")
// 	}

// 	if token := b.conn.Publish(topic.String(), b.qos, false, bytes); token.Wait() && token.Error() != nil {
// 		return token.Error()
// 	}
// 	return nil
// }
