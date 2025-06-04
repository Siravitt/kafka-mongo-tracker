package kafka

import (
	"log/slog"

	"github.com/IBM/sarama"
	"github.com/Siravitt/kafka-mongo-tracker/config"
)

// var config config.Config

func NewConsumerGroupWithConfig(group string, addrs config.KafKa) (sarama.ConsumerGroup, error) {
	cfg := sarama.NewConfig()

	cfg.Version = sarama.V2_8_0_0

	// consumer group
	cfg.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	cfg.Consumer.Group.Session.Timeout = consumerGroupSessionTimeout
	cfg.Consumer.Group.Heartbeat.Interval = consumerGroupHeartbeatInterval

	// message
	cfg.Consumer.Fetch.Min = consumerFetchMin
	cfg.Consumer.Fetch.Default = consumerFetchDefault

	// offset
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	cfg.Consumer.Offsets.AutoCommit.Enable = consumerOffsetsAutoCommitEnable
	cfg.Consumer.Offsets.AutoCommit.Interval = consumerOffsetsAutoCommitInterval

	// retry
	cfg.Consumer.Retry.Backoff = consumerRetryBackoff

	consumerGroup, err := sarama.NewConsumerGroup(addrs.ConsumerURL, group, cfg)
	if err != nil {
		slog.Error("failed to create Kafka consumer group",
			"error", err,
			"group", group,
			"broker", addrs.ConsumerURL,
		)
		return nil, err
	}

	return consumerGroup, nil
}
