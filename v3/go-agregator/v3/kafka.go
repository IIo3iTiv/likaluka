package main

import (
	"log"

	"github.com/IBM/sarama"
)

type DataJSON struct {
	Guid uint32      `json:"guid"`
	Data []BatchJSON `json:"data"`
}

type BatchJSON struct {
	Date int64   `json:"date"`
	R    float32 `json:"r"`
	S    float32 `json:"s"`
	T    float32 `json:"t"`
}

const (
	kafka0    string = "192.168.1.114:9092"
	kafka1    string = "192.168.1.154:9092"
	kafka2    string = "192.168.1.51:9092"
	topicName string = "itcamp"
)

var k Kafka

type Kafka struct {
	producer sarama.SyncProducer
}

func InitKafka() (err error) {
	brokers := []string{kafka0, kafka1, kafka2}
	config := sarama.NewConfig()
	config.Net.SASL.Enable = false
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	k.producer, err = sarama.NewSyncProducer(brokers, config)
	return
}

func KafkaPublish(message string) (err error) {
	msg := &sarama.ProducerMessage{
		Topic: topicName,
		Value: sarama.StringEncoder(message),
	}
	_, _, err = k.producer.SendMessage(msg)
	if err != nil {
		log.Printf("Error: Kafka: Publish: %s", err.Error())
		return
	}
	return
}
