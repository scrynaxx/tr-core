package events

import "time"

// Config задаёт параметры доставки событий и работы с RabbitMQ.
// Нулевые значения отдельных параметров заменяются значениями по умолчанию.
type Config struct {
	PublisherPoolSize    int
	ConnectionRetryDelay time.Duration
	ConnectionTimeout    time.Duration
	Outbox               OutboxConfig
	Subscriber           SubscriberConfig
}

// OutboxConfig задаёт параметры аренды и публикации outbox-событий.
type OutboxConfig struct {
	BatchSize       int
	Concurrency     int
	PollInterval    time.Duration
	PublishTimeout  time.Duration
	RetryDelay      time.Duration
	MaxRetries      int
	LeaseMargin     time.Duration
	Retention       time.Duration
	CleanupInterval time.Duration
}

// SubscriberConfig задаёт retry policy входящих событий.
type SubscriberConfig struct {
	RetryDelay        time.Duration
	MaxAttempts       int
	DeadLetterTimeout time.Duration
}

func defaultConfig() Config {
	return Config{
		PublisherPoolSize:    5,
		ConnectionRetryDelay: 100 * time.Millisecond,
		ConnectionTimeout:    10 * time.Second,
		Outbox: OutboxConfig{
			BatchSize:       50,
			Concurrency:     5,
			PollInterval:    500 * time.Millisecond,
			PublishTimeout:  15 * time.Second,
			RetryDelay:      5 * time.Second,
			MaxRetries:      5,
			LeaseMargin:     30 * time.Second,
			Retention:       7 * 24 * time.Hour,
			CleanupInterval: time.Hour,
		},
		Subscriber: SubscriberConfig{
			RetryDelay:        15 * time.Second,
			MaxAttempts:       5,
			DeadLetterTimeout: 15 * time.Second,
		},
	}
}

func normalizeConfig(config *Config) Config {
	result := defaultConfig()
	if config == nil {
		return result
	}

	if config.PublisherPoolSize > 0 {
		result.PublisherPoolSize = config.PublisherPoolSize
	}
	if config.ConnectionRetryDelay > 0 {
		result.ConnectionRetryDelay = config.ConnectionRetryDelay
	}
	if config.ConnectionTimeout > 0 {
		result.ConnectionTimeout = config.ConnectionTimeout
	}
	if config.Outbox.BatchSize > 0 {
		result.Outbox.BatchSize = config.Outbox.BatchSize
	}
	if config.Outbox.Concurrency > 0 {
		result.Outbox.Concurrency = config.Outbox.Concurrency
	}
	if config.Outbox.PollInterval > 0 {
		result.Outbox.PollInterval = config.Outbox.PollInterval
	}
	if config.Outbox.PublishTimeout > 0 {
		result.Outbox.PublishTimeout = config.Outbox.PublishTimeout
	}
	if config.Outbox.RetryDelay > 0 {
		result.Outbox.RetryDelay = config.Outbox.RetryDelay
	}
	if config.Outbox.MaxRetries > 0 {
		result.Outbox.MaxRetries = config.Outbox.MaxRetries
	}
	if config.Outbox.LeaseMargin > 0 {
		result.Outbox.LeaseMargin = config.Outbox.LeaseMargin
	}
	if config.Outbox.Retention > 0 {
		result.Outbox.Retention = config.Outbox.Retention
	}
	if config.Outbox.CleanupInterval > 0 {
		result.Outbox.CleanupInterval = config.Outbox.CleanupInterval
	}
	if config.Subscriber.RetryDelay > 0 {
		result.Subscriber.RetryDelay = config.Subscriber.RetryDelay
	}
	if config.Subscriber.MaxAttempts > 0 {
		result.Subscriber.MaxAttempts = config.Subscriber.MaxAttempts
	}
	if config.Subscriber.DeadLetterTimeout > 0 {
		result.Subscriber.DeadLetterTimeout = config.Subscriber.DeadLetterTimeout
	}

	return result
}
