package kafka

import "time"

const (
	consumerGroupSessionTimeout       = 10 * time.Second
	consumerGroupHeartbeatInterval    = 3 * time.Second
	consumerOffsetsAutoCommitEnable   = true
	consumerOffsetsAutoCommitInterval = 1 * time.Second
	consumerRetryBackoff              = 500 * time.Millisecond
	consumerFetchMin                  = 1
	consumerFetchDefault              = 1 * 1024 * 1024

	producerReturnSuccesses = true
	producerRetryMax        = 5
)
