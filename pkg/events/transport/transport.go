package transport

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	// ErrUnavailable сообщает о временной недоступности AMQP-соединения.
	ErrUnavailable = errors.New("event bus unavailable")

	// ErrUnroutable сообщает, что RabbitMQ не нашёл очередь для routing key.
	ErrUnroutable = errors.New("event has no route")
)

// Config содержит параметры AMQP-подключения.
type Config struct {
	User                 string
	Password             string
	Address              string
	Vhost                string
	PublisherPoolSize    int
	ConnectionRetryDelay time.Duration
	ConnectionTimeout    time.Duration
}

// Transport управляет AMQP-соединением и подтверждаемой публикацией сообщений.
type Transport struct {
	mu         sync.Mutex
	connection *amqp.Connection

	publishers    []*publisher
	publisherPool chan *publisher

	url                  string
	config               amqp.Config
	connectionRetryDelay time.Duration
	connectionTimeout    time.Duration
}

// New открывает AMQP-соединение и создаёт пул publisher-каналов.
func New(ctx context.Context, config Config) (*Transport, error) {
	connectionURL := (&url.URL{
		Scheme: "amqp",
		User:   url.UserPassword(config.User, config.Password),
		Host:   config.Address,
	}).String()
	amqpConfig := amqp.Config{Vhost: config.Vhost}

	transport := &Transport{
		publishers:           make([]*publisher, config.PublisherPoolSize),
		publisherPool:        make(chan *publisher, config.PublisherPoolSize),
		url:                  connectionURL,
		config:               amqpConfig,
		connectionRetryDelay: config.ConnectionRetryDelay,
		connectionTimeout:    config.ConnectionTimeout,
	}
	for i := range transport.publishers {
		transport.publishers[i] = &publisher{}
		transport.publisherPool <- transport.publishers[i]
	}

	connection, err := transport.dial(ctx)
	if err != nil {
		return nil, err
	}

	transport.connection = connection
	return transport, nil
}

// Channel открывает AMQP-канал, ожидая восстановления соединения до отмены операции.
func (c *Transport) Channel(ctx context.Context) (*amqp.Channel, error) {
	for ctx.Err() == nil {
		connection, err := c.getConnection(ctx)
		if err == nil {
			channel, channelErr := connection.Channel()
			if channelErr == nil {
				return channel, nil
			}
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(c.connectionRetryDelay):
		}
	}

	return nil, ctx.Err()
}

// Publish отправляет persistent-сообщение и ожидает publisher confirm от RabbitMQ.
func (c *Transport) Publish(ctx context.Context, exchange, routingKey string, body []byte) (err error) {
	var selected *publisher
	select {
	case <-ctx.Done():
		return fmt.Errorf("%w: %w", ErrUnavailable, ctx.Err())
	case selected = <-c.publisherPool:
	}
	defer func() {
		c.publisherPool <- selected
	}()

	connection, err := c.getConnection(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrUnavailable, err)
	}

	return selected.publish(ctx, connection, exchange, routingKey, body)
}

// Close закрывает publisher, затем AMQP-соединение.
func (c *Transport) Close() {
	for i := range c.publishers {
		c.publishers[i].close()
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.connection != nil && !c.connection.IsClosed() {
		_ = c.connection.Close()
	}
}

func (c *Transport) getConnection(ctx context.Context) (*amqp.Connection, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connection != nil && !c.connection.IsClosed() {
		return c.connection, nil
	}

	connection, err := c.dial(ctx)
	if err != nil {
		return nil, fmt.Errorf("dial rabbitmq: %w", err)
	}

	c.connection = connection
	return connection, nil
}

func (c *Transport) dial(ctx context.Context) (*amqp.Connection, error) {
	dialCtx, cancel := context.WithTimeout(ctx, c.connectionTimeout)
	defer cancel()

	config := c.config
	config.Dial = func(network, address string) (net.Conn, error) {
		connection, err := (&net.Dialer{}).DialContext(dialCtx, network, address)
		if err != nil {
			return nil, err
		}
		if err = connection.SetDeadline(time.Now().Add(c.connectionTimeout)); err != nil {
			_ = connection.Close()
			return nil, err
		}
		return connection, nil
	}

	return amqp.DialConfig(c.url, config)
}
