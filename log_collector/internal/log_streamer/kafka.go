package log_streamer

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/segmentio/kafka-go"
)

type KafkaLogProducer struct {
    topicName  string
    kafkaWriter *kafka.Writer
}

func createNewKafkaTopic(broker string,topic string)error{

    conn, err := kafka.Dial("tcp", broker)
    if err != nil {
        return fmt.Errorf("failed to connect to Kafka broker: %v", err)
    }
    defer conn.Close()

    controller, err := conn.Controller()
    if err != nil {
        return fmt.Errorf("failed to get controller: %v", err)
    }

    controllerConn, err := kafka.Dial("tcp", controller.Host)
    if err != nil {
        return fmt.Errorf("failed to connect to controller: %v", err)
    }
    defer controllerConn.Close()

    topicConfigs := []kafka.TopicConfig{
        {
            Topic:             topic,
            NumPartitions:     3,
            ReplicationFactor: 1,
        },
    }

    err = controllerConn.CreateTopics(topicConfigs...)
    if err != nil {
        return fmt.Errorf("failed to create topic: %v", err)
    }

    fmt.Println("✅ Topic created successfully:", topic)
    return nil
}



// SendLog sends a message to the Kafka topic
func (p *KafkaLogProducer) SendLog(msg string) error {
    err := p.kafkaWriter.WriteMessages(context.Background(),
        kafka.Message{Value: []byte(msg)},
    )
    if err != nil {
        return fmt.Errorf("failed to send message: %v", err)
    }
    fmt.Println("Message delivered successfully")
    return nil
}

func (p *KafkaLogProducer)Close()error{
    err:=p.kafkaWriter.Close()
    return err
}

// NewKafkaLogProducer initializes and returns a KafkaLogProducer
func NewKafkaLogProducer(topicName string) (*KafkaLogProducer,error) {
    broker := os.Getenv("KAFKA_BROKER")
    if broker == "" {
        return nil,fmt.Errorf("kafka broker url not found")
    }

    err:=createNewKafkaTopic(broker,topicName)
    if err !=nil {
        return nil,fmt.Errorf("error in creating kafka topic:%v",err)
    }

    writer := kafka.NewWriter(kafka.WriterConfig{
        Brokers:  []string{broker},
        Topic:    topicName,
        Balancer: &kafka.LeastBytes{},
    })

    log.Printf("Kafka producer initialized for topic: %s\n", topicName)

    return &KafkaLogProducer{
        topicName:  topicName,
        kafkaWriter: writer,
    },nil
}
