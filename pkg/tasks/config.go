package tasks

import "time"

// Config определяет ограничения, lease и политику повторов task runner-а.
type Config struct {
	// Количество задач, получаемых из PostgreSQL одним claim-запросом.
	BatchSize int

	// Общее количество одновременно выполняемых задач на одной реплике.
	Concurrency int

	// Интервал поиска готовых задач.
	PollInterval time.Duration

	// Срок, в течение которого runner обязан продлить lease выполняемой задачи.
	LeaseDuration time.Duration

	// Базовая задержка экспоненциального backoff после ошибки.
	RetryDelay time.Duration

	// Количество попыток до окончательного перевода задачи в failed.
	MaxAttempts int

	// Срок хранения завершённых и окончательно упавших задач.
	Retention time.Duration

	// Период удаления задач, срок хранения которых истёк.
	CleanupInterval time.Duration
}

func defaultConfig() Config {
	return Config{
		BatchSize:       300,
		Concurrency:     300,
		PollInterval:    100 * time.Millisecond,
		LeaseDuration:   3 * time.Minute,
		RetryDelay:      5 * time.Second,
		MaxAttempts:     10,
		Retention:       72 * time.Hour,
		CleanupInterval: time.Hour,
	}
}

func normalizeConfig(config *Config) Config {
	result := defaultConfig()
	if config == nil {
		return result
	}

	if config.BatchSize > 0 {
		result.BatchSize = config.BatchSize
	}
	if config.Concurrency > 0 {
		result.Concurrency = config.Concurrency
	}
	if config.PollInterval > 0 {
		result.PollInterval = config.PollInterval
	}
	if config.LeaseDuration > 0 {
		result.LeaseDuration = config.LeaseDuration
	}
	if config.RetryDelay > 0 {
		result.RetryDelay = config.RetryDelay
	}
	if config.MaxAttempts > 0 {
		result.MaxAttempts = config.MaxAttempts
	}
	if config.Retention > 0 {
		result.Retention = config.Retention
	}
	if config.CleanupInterval > 0 {
		result.CleanupInterval = config.CleanupInterval
	}

	return result
}
