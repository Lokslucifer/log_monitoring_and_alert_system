package alertdispatcher

import (
	"fmt"
	"os"
	"log_processor/internal/models"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQAlertPublisher struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	qname string
}

// PublishLog sends a log entry to RabbitMQ
func (r *RabbitMQAlertPublisher) PublishLog(logEntry models.LogEntry) error {

	logMsg := fmt.Sprintf(
		"Timestamp:%v - Level:%v - Message:%v - SourceFile:%v - LineNumber:%v",
		logEntry.Timestamp,
		logEntry.Level,
		logEntry.Message,
		logEntry.SourceFile,
		logEntry.LineNumber,
	)

	err := r.ch.Publish(
		"",        // exchange
		r.qname,   // routing key (queue name)
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(logMsg),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish log: %v", err)
	}

	return nil
}


func (r *RabbitMQAlertPublisher)Close(){
	r.ch.Close()
	r.conn.Close()
}
// NewRabbitMQAlertPublisher initializes and returns a RabbitMQ publisher
func NewRabbitMQAlertPublisher(queueName string) (AlertPublisher, error) {
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

	return &RabbitMQAlertPublisher{
		conn:  conn,
		ch:    ch,
		qname: q.Name,
	}, nil
}
