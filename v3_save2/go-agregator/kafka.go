package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/IBM/sarama"
)

const (
	kafka0    string = "kafka:9091"
	kafka1    string = "kafka-1:9092"
	kafka2    string = "kafka-2:9092"
	topicName string = "itcamp"
	user      string = "sa"
	pass      string = "000000"
)

type Kafka struct {
	producer sarama.SyncProducer
}

func InitKafka() (k Kafka, err error) {
	brokers := []string{kafka0 /*, kafka1, kafka2*/}
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Net.SASL.Enable = true
	config.Net.SASL.User = user
	config.Net.SASL.Password = pass
	config.Net.SASL.Mechanism = sarama.SASLTypePlaintext

	certs := x509.NewCertPool()
	pemPath := filepath.Join(GetExecuteFilePath(), "kafka.truststore.pem")
	pemData, err := os.ReadFile(pemPath)
	if err != nil {
		fmt.Printf("Couldn't load cert: %s\n", err.Error())
		return
	}
	certs.AppendCertsFromPEM(pemData)

	config.Net.TLS.Enable = true
	config.Net.TLS.Config = &tls.Config{
		RootCAs:            certs,
		InsecureSkipVerify: true,
	}

	k.producer, err = sarama.NewSyncProducer(brokers, config)
	return
}

func (k *Kafka) Publish(message string) {
	msg := &sarama.ProducerMessage{
		Topic: topicName,
		Value: sarama.StringEncoder(message),
	}
	p, o, err := k.producer.SendMessage(msg)
	if err != nil {
		log.Printf("Error: Kafka: Publish: %s", err.Error())
		return
	}
	log.Printf("Partition: %d, Offset: %d\n", p, o)
}
