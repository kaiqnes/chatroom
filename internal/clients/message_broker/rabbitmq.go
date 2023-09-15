package message_broker

import (
	"chatroom/internal/config"
	"chatroom/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
)

type mq struct {
	channel *amqp.Channel
	queue   amqp.Queue
	socket  *socketio.Server
}

func (m *mq) Send(message []byte) error {
	err := m.channel.PublishWithContext(context.Background(),
		"",
		m.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})
	if err != nil {
		return err
	}
	return nil
}

func (m *mq) Listen() error {
	msgs, err := m.channel.Consume(
		m.queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Create a channel to receive OS signals
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for msg := range msgs {
			fmt.Printf("Received message: %s\n", msg.Body)
			var botMsg domain.MessageRequestDto
			err = json.Unmarshal(msg.Body, &botMsg)
			if err != nil {
				fmt.Printf("error to unmarshal message: %+v\n", err)
				continue
			}
			m.socket.OnConnect("/", func(s socketio.Conn) error {
				fmt.Printf("msg to be sent: %+v\n", botMsg)
				s.SetContext("")
				b := m.socket.BroadcastToNamespace("/", "chat message", botMsg)
				fmt.Printf("broadcast: %+v\n", b)
				return nil
			})
		}
	}()

	// Block the execution until a termination signal is received
	<-sig
	return nil
}

func NewMQ(cfg *config.Config, socketServer *socketio.Server) (domain.MessageQueue, error) {
	host := fmt.Sprintf(cfg.RabbitMQHostTemplate, cfg.RabbitMQUser, cfg.RabbitMQPassword, cfg.RabbitMQHost, cfg.RabbitMQPort)

	conn, err := amqp.Dial(host)
	if err != nil {
		return nil, err
	}
	//defer conn.Close() // Keep it open

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	//defer channel.Close() // Keep it open

	queue, err := channel.QueueDeclare(
		cfg.RabbitMQQueue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &mq{
		channel: channel,
		queue:   queue,
		socket:  socketServer,
	}, nil
}
