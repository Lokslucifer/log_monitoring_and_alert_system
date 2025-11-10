package alertconsumer

import (
	"fmt"
	"os"
	"alert_system/internal/alert_sender"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQAlertConsumer struct {
	conn  *amqp.Connection
	alerts   <-chan amqp.Delivery
	ch *amqp.Channel
	qname string
}

// PublishLog sends a log entry to RabbitMQ
func (r *RabbitMQAlertConsumer) StartConsumingLog(alerter alertsender.AlertSender) {

		for alert := range r.alerts {
			alerter.SendAlert(string(alert.Body))
		}

}


func (r *RabbitMQAlertConsumer)Close(){
	r.ch.Close()
	r.conn.Close()

}
// NewRabbitMQAlertConsumer initializes and returns a RabbitMQ publisher
func NewRabbitMQAlertConsumer(queueName string) (*RabbitMQAlertConsumer, error) {
	broker := os.Getenv("RABBIT_MQ_BROKER")
	user := os.Getenv("RABBIT_MQ_USER")
	password := os.Getenv("RABBIT_MQ_PASSWORD")

	if user == "" {
		log.Fatal("RABBIT_MQ_USER is not set in environment")
	}

	if password == "" {
		log.Fatal("RABBIT_MQ_PASSWORD is not set in environment")
	}

	if broker == "" {
		log.Fatal("RABBIT_MQ_BROKER is not set in environment")
	}

	// Connect to RabbitMQ
	url := fmt.Sprintf("amqp://%s:%s@%s/", user, password, broker)
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %v", err)
	}

	// Declare queue
	q, err := ch.QueueDeclare(
		queueName, // queue name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %v", err)
	}

		// Consume messages
	msgs, err := ch.Consume(
		q.Name,
		"",    // consumer tag
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {

		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to register consumer: %v", err)
	}

	return &RabbitMQAlertConsumer{
		conn:  conn,
		ch:    ch,
		qname: q.Name,
		alerts: msgs,

	}, nil
}
