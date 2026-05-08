package rabbitmq

import (
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ResilientClient struct {
	url           string
	Connection    *amqp.Connection
	connCloseChan chan *amqp.Error
	hooks         []func(*amqp.Connection)
}

func NewResilientClient(url string) *ResilientClient {
	c := &ResilientClient{
		url: url,
	}
	c.connect()
	go c.reconnectLoop()
	return c
}

func (c *ResilientClient) RegisterReconnectHook(hook func(*amqp.Connection)) {
	c.hooks = append(c.hooks, hook)
	// If already connected, trigger immediately for the first setup
	if c.Connection != nil && !c.Connection.IsClosed() {
		go hook(c.Connection)
	}
}

func (c *ResilientClient) connect() {
	for {
		conn, err := amqp.Dial(c.url)
		if err == nil {
			c.Connection = conn
			c.connCloseChan = make(chan *amqp.Error)
			c.Connection.NotifyClose(c.connCloseChan)
			slog.Info("Successfully connected to RabbitMQ")
			return
		}
		slog.Error("Failed to connect to RabbitMQ, retrying in 5 seconds", "error", err)
		time.Sleep(5 * time.Second)
	}
}

func (c *ResilientClient) reconnectLoop() {
	for {
		err, ok := <-c.connCloseChan
		if !ok {
			slog.Info("RabbitMQ connection close channel closed")
			return
		}
		if err != nil {
			slog.Warn("RabbitMQ connection closed unexpectedly, attempting to reconnect", "error", err)
			c.connect()
			slog.Info("RabbitMQ reconnected. Triggering reconnect hooks (Auto-Redeclaration)")
			for _, hook := range c.hooks {
				go hook(c.Connection) // Execute hooks asynchronously
			}
		} else {
			slog.Info("RabbitMQ connection closed gracefully")
			return
		}
	}
}

func (c *ResilientClient) Close() {
	if c.Connection != nil && !c.Connection.IsClosed() {
		c.Connection.Close()
	}
}
