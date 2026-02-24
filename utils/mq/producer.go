package mq

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type Producer struct {
	sarama.AsyncProducer
}

func NewProducer() *Producer {
	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForLocal // use Boolean filter instead
	conf.Producer.Return.Errors = true
	conf.Producer.Return.Successes = true
	conf.Producer.Compression = sarama.CompressionZSTD
	conf.Producer.Flush.Messages = 100
	conf.Producer.Flush.Frequency = 10 * time.Millisecond
	conf.Producer.Flush.Bytes = 64 * 1024
	conf.Producer.Retry.Backoff = 500 * time.Millisecond
	producer, err := sarama.NewAsyncProducer([]string{"192.168.0.100:9092"}, conf)
	if errors.Is(err, sarama.ErrClosedClient) {
		panic(errors.New("error in NewAsyncProducer :" + err.Error()))
	}
	pd := &Producer{AsyncProducer: producer}
	pd.handleResponses()
	return pd
}

func (p *Producer) handleResponses() {
	go func() {
		for {
			select {
			case err := <-p.AsyncProducer.Errors():
				log.Printf("Kafka produce handle response error%v\n", err)
			case <-p.AsyncProducer.Successes():
			}
		}
	}()
}

func (p *Producer) Send(topic string, value any) {
	v, err := json.Marshal(value)
	if err != nil {
		log.Printf("Json marshal error %v\n", err)
	}

	p.AsyncProducer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(v),
	}
}
