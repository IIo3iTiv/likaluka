package producer

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"kafka_yuacloud/scram"
	"os"
	"path/filepath"

	"github.com/IBM/sarama"
)

var brokers []string = []string{
	"rc1a-4kmc3iak927gt9nm.mdb.yandexcloud.net:9091",
	"rc1b-mcifivrkhj6d6vd3.mdb.yandexcloud.net:9091",
	"rc1d-53oqpr9s7ge335kh.mdb.yandexcloud.net:9091",
}

const topic = "test_case"
const u = "oleg"
const p = "ECsFs9%WSLmDtK5%PEnD^U3L3??~7d-"

func Start() error {
	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Return.Successes = true
	conf.Version = sarama.V0_10_0_0
	conf.ClientID = "sasl_scram_client"
	conf.Net.SASL.Enable = true
	conf.Net.SASL.Handshake = true
	conf.Net.SASL.User = u
	conf.Net.SASL.Password = p
	conf.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &scram.XDGSCRAMClient{HashGeneratorFcn: scram.SHA512} }
	conf.Net.SASL.Mechanism = sarama.SASLMechanism(sarama.SASLTypeSCRAMSHA512)

	certs := x509.NewCertPool()
	pemPath := filepath.Join(getExecuteFilePath(), "..", "cert", "YandexInternalRootCA.crt")
	pemData, err := os.ReadFile(pemPath)
	if err != nil {
		fmt.Printf("Couldn't load cert: %s\n", err.Error())
		return err
	}
	certs.AppendCertsFromPEM(pemData)

	conf.Net.TLS.Enable = true
	conf.Net.TLS.Config = &tls.Config{
		RootCAs:            certs,
		InsecureSkipVerify: true,
	}

	syncProducer, err := sarama.NewSyncProducer(brokers, conf)
	if err != nil {
		fmt.Printf("Couldn't create producer: %s\n", err.Error())
		return err
	}
	publish(topic, "test_message", syncProducer)

	return nil
}

func StartDS() error {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Ждем подтверждения от всех реплик
	config.Producer.Retry.Max = 5                    // Количество попыток
	config.Producer.Return.Successes = true

	// Если используется SASL аутентификация (часто в Yandex Cloud)
	config.Net.SASL.Enable = true
	config.Net.SASL.User = u
	config.Net.SASL.Password = p
	config.Net.SASL.Mechanism = sarama.SASLTypePlaintext

	certs := x509.NewCertPool()
	pemPath := filepath.Join(getExecuteFilePath(), "..", "cert", "YandexInternalRootCA.crt")
	pemData, err := os.ReadFile(pemPath)
	if err != nil {
		fmt.Printf("Couldn't load cert: %s\n", err.Error())
		return err
	}
	certs.AppendCertsFromPEM(pemData)

	// Если используется SSL
	config.Net.TLS.Enable = true
	config.Net.TLS.Config = &tls.Config{
		RootCAs:            certs,
		InsecureSkipVerify: true,
	}

	// Создаем продюсера
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		fmt.Printf("Ошибка создания продюсера: %v", err)
		return err
	}
	defer producer.Close()

	publish(topic, "test_message", producer)
	return nil
}

func publish(topic, message string, producer sarama.SyncProducer) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	p, o, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Printf("Error publish: %s\n", err.Error())
		return err
	}

	fmt.Printf("Partition: %d, Offset: %d\n", p, o)
	return nil
}

func getExecuteFilePath() string {
	ex, _ := os.Executable()
	return filepath.Dir(ex)
}
