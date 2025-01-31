package queue

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
)

type RabbitMQ struct {
	Logger *logrus.Logger
}

func (r *RabbitMQ) newConnection() *amqp.Connection {
	host := os.Getenv("RABBITMQ_HOST")
	conn, err := amqp.Dial(host)
	r.failOnError(err, "Failed to connect to RabbitMQ")

	return conn
}

func (r *RabbitMQ) CreateQueue(name string) error {
	conn:= r.newConnection()
	ch, err := conn.Channel()
	r.failOnError(err, "Failed to open a channel")
	defer ch.Close()
	defer conn.Close()

	_, err = ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	r.failOnError(err, "Failed to declare a queue")

	return err
}
func (r *RabbitMQ) PublishMessage(queueName string, bodyJson string) {
	conn:= r.newConnection()
	ch, err := conn.Channel()
	r.failOnError(err, "Failed to open a channel")
	defer ch.Close()
	defer conn.Close()

	err = ch.Publish(
		"",         // exchange
		queueName,  // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(bodyJson),
		})
	r.failOnError(err, "Failed to publish a message")
}

func (r *RabbitMQ) failOnError(err error, msg string) {
	if err != nil {
		r.Logger.Fatalf("%s: %s", msg, err)
	}
}
