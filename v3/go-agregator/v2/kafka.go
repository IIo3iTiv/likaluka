package main

import (
	"log"

	"github.com/IBM/sarama"
)

const (
	kafka0    string = "192.168.1.114:9092"
	kafka1    string = "192.168.1.154:9092"
	kafka2    string = "192.168.1.51:9092"
	topicName string = "itcamp"
)

type Kafka struct {
	producer sarama.SyncProducer
}

func InitKafka() (k Kafka, err error) {
	brokers := []string{kafka0, kafka1, kafka2}
	config := sarama.NewConfig()
	config.Net.SASL.Enable = false
	// config.Net.SASL.User = user
	// config.Net.SASL.Password = pass
	// config.Net.SASL.Mechanism = sarama.Plain

	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Partitioner = sarama.NewRandomPartitioner

	// certs := x509.NewCertPool()
	// pemPath := filepath.Join(GetExecuteFilePath(), "kafka.truststore.pem")
	// pemData, err := os.ReadFile(pemPath)
	// if err != nil {
	// 	fmt.Printf("Couldn't load cert: %s\n", err.Error())
	// 	return
	// }
	// certs.AppendCertsFromPEM(pemData)

	// config.Net.TLS.Enable = true
	// config.Net.TLS.Config = &tls.Config{
	// RootCAs:            certs,
	// InsecureSkipVerify: true,
	// }

	k.producer, err = sarama.NewSyncProducer(brokers, config)
	return
}

func (k *Kafka) Publish(message string) {
	msg := &sarama.ProducerMessage{
		Topic: topicName,
		Value: sarama.StringEncoder(message),
	}
	_, _, err := k.producer.SendMessage(msg)
	if err != nil {
		log.Printf("Error: Kafka: Publish: %s", err.Error())
		return
	}
	// log.Printf("Partition: %d, Offset: %d\n", p, o)
}
