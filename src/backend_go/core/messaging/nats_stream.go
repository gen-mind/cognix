package messaging

import (
	"context"
	proto2 "github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	_ "github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
	"sync"
	"time"
)

// clientStream is a struct that represents a client stream.
// It provides methods to interact with the stream, such as publishing messages,
// listening for messages, and managing the stream configuration and status.
type clientStream struct {
	js        jetstream.JetStream
	conn      *nats.Conn
	cancel    context.CancelFunc
	ctx       context.Context
	wg        *sync.WaitGroup
	streamCfg *StreamConfig
	ackWait   time.Duration
}

// IsOnline checks if the client is currently online by checking the connection status.
func (c *clientStream) IsOnline() bool {
	return c.conn.Status() == nats.CONNECTED
}

// StreamConfig returns the StreamConfig object associated with the clientStream
func (c *clientStream) StreamConfig() *StreamConfig {

	return c.streamCfg
}

// Close completes the current operation and waits for all goroutines to finish.
// It increments the wait group by 1, cancels the context, and waits for all goroutines to finish.
// Any goroutine that is waiting for the wait group to become zero will unblock after this operation completes.
func (c *clientStream) Close() {
	c.wg.Add(1)
	c.cancel()
	c.wg.Wait()
}

// Publish sends a message to the specified stream and topic.
// It first creates or updates the stream using the provided stream name,
// with a retention policy of WorkQueuePolicy and the specified topic.
// Then it marshals the provided body into a protobuf message,
// and publishes it to the NATS JetStream using the specified context and topic.
// Returns any error encountered during the process.
func (c *clientStream) Publish(ctx context.Context, streamName, topic string, body proto2.Message) error {
	_, err := c.js.CreateOrUpdateStream(context.Background(), jetstream.StreamConfig{
		Name:      streamName,
		Retention: jetstream.WorkQueuePolicy,
		Subjects:  []string{topic},
	})
	if err != nil {
		zap.S().Errorf("Error creating stream: %s", err.Error())
		return err
	}

	message, err := proto2.Marshal(body)
	if err != nil {
		return err
	}
	_, err = c.js.Publish(ctx, topic, message)
	if err != nil {
		return err
	}
	return nil
}

func (c *clientStream) Listen(ctx context.Context, streamName, topic string, handler MessageHandler) error {

	stream, err := c.js.CreateOrUpdateStream(context.Background(), jetstream.StreamConfig{
		Name:      streamName,
		Retention: jetstream.WorkQueuePolicy,
		Subjects:  []string{topic},
	})
	if err != nil {
		zap.S().Errorf("Error creating stream: %s", err.Error())
		return err
	}

	cons, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       streamName,
		MaxDeliver:    c.streamCfg.MaxDeliver,
		FilterSubject: topic,
		AckPolicy:     jetstream.AckExplicitPolicy,
		AckWait:       c.ackWait,
		DeliverPolicy: jetstream.DeliverAllPolicy,
	})
	if err != nil {
		zap.S().Errorf("Failed to create consumer for subscription %v", err)
	}
	cons.Consume(func(msg jetstream.Msg) {
		msg.InProgress()
		if err := handler(ctx, msg); err != nil {
			zap.S().Errorf("Error handling message: %s", err.Error())
		}
		err := msg.Ack()
		if err != nil {
			zap.S().Errorf("Error acknowledging message: %s", err.Error())
		}

	})
	<-c.ctx.Done()
	c.wg.Done()
	return nil
}

// NewClientStream creates a new client stream by connecting to NATS using the provided configuration.
// It returns a Client interface and an error. If the connection to NATS fails, it returns nil for the client
// and the corresponding error. If the connection is successful, it creates a JetStream object and initializes
// a clientStream object with the JetStream and other necessary attributes.
// The clientStream object implements the Client interface and is returned along with a nil error.
//
// Example usage:
//
//	cfg := &Config{
//	    Nats: &natsConfig{
//	        URL:                 "localhost:4222",
//	        ConnectorStreamName: "test-1",
//	    },
//	    Stream: &StreamConfig{},
//	}
//
// client, err := NewClientStream(cfg)
//
//	if err != nil {
//	    // Handle error
//	}
//
// defer client.Close()
func NewClientStream(cfg *Config) (Client, error) {
	conn, err := nats.Connect(
		cfg.Nats.URL,
		nats.Timeout(time.Duration(cfg.Nats.ConnectTimeout)*time.Second),
		nats.ReconnectWait(time.Duration(cfg.Nats.ReconnectTimeout)*time.Second),
		nats.MaxReconnects(cfg.Nats.MaxReconnectAttempts),
	)
	if err != nil {
		zap.S().Errorf("Error connecting to NATS: %s", err.Error())
		return nil, err
	}
	js, err := jetstream.New(conn)
	if err != nil {
		zap.S().Errorf("Error connecting to NATS: %s", err.Error())
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &clientStream{
		js:        js,
		ctx:       ctx,
		conn:      conn,
		cancel:    cancel,
		wg:        &sync.WaitGroup{},
		streamCfg: cfg.Stream,
		ackWait:   time.Duration(cfg.Stream.AckWait) * time.Second,
	}, nil
}
