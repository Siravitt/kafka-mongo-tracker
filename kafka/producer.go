package kafka

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/Siravitt/kafka-mongo-tracker/config"
)

func NewSyncProducerGuarantee(addrs config.KafKa) sarama.SyncProducer {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	cfg.Producer.Return.Successes = producerReturnSuccesses
	cfg.Producer.Retry.Max = producerRetryMax

	producer, err := sarama.NewSyncProducer(addrs.ProducerURL, cfg)
	if err != nil {
		log.Panic(err)
	}
	return producer
}
