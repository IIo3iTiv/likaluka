package consumer

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"kafka_yuacloud/scram"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/IBM/sarama"
)

var brokers []string = []string{
	"rc1a-4kmc3iak927gt9nm.mdb.yandexcloud.net:9091",
	"rc1b-mcifivrkhj6d6vd3.mdb.yandexcloud.net:9091",
	"rc1d-53oqpr9s7ge335kh.mdb.yandexcloud.net:9091",
}

func Start() error {
	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Return.Errors = true
	conf.Version = sarama.V0_10_0_0
	conf.ClientID = "sasl_scram_client"
	conf.Metadata.Full = true
	conf.Net.SASL.Enable = true
	conf.Net.SASL.Handshake = true
	conf.Net.SASL.User = "oleg"
	conf.Net.SASL.Password = "ECsFs9%WSLmDtK5%PEnD^U3L3??~7d-"
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
		RootCAs: certs,
	}

	master, err := sarama.NewConsumer(brokers, conf)
	if err != nil {
		fmt.Printf("Couldn't create consumer: %s\n", err.Error())
		return err
	}
	defer func() {
		if err := master.Close(); err != nil {
			panic(err)
		}
	}()

	topic := "test_case"

	consumer, err := master.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return err
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	msgCount := 0
	doneCh := make(chan struct{})

	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Printf("ERROR: %s\n", err.Error())
			case msg := <-consumer.Messages():
				msgCount++
				fmt.Printf("Recieved messages Key: {%s}, Value: {%s}\n", string(msg.Key), string(msg.Value))
			case <-signals:
				fmt.Println("Interrupt is detached")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	fmt.Printf("Proceddes %d messages\n", msgCount)
	return nil
}

func getExecuteFilePath() string {
	ex, _ := os.Executable()
	return filepath.Dir(ex)
}
