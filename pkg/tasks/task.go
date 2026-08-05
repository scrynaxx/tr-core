package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// ErrLeaseLost означает, что попытка выполнения больше не владеет lease задачи.
var ErrLeaseLost = errors.New("task lease lost")

// Task описывает задачу, которую runner должен выполнить не раньше заданного времени.
type Task struct {
	// Тип задачи и зарегистрированный для неё handler.
	Source string `json:"source"`

	// Уникальный идентификатор задачи внутри source для идемпотентной постановки.
	SourceID string `json:"source_id"`

	// Ключ группы строго последовательного выполнения; nil не ограничивает параллелизм.
	SerialKey *string `json:"serial_key,omitempty"`

	// Входные данные handler-а в формате JSON; nil используется для задач без входных данных.
	Payload json.RawMessage `json:"payload"`

	// Самое раннее время начала выполнения.
	RunAt time.Time `json:"run_at"`
}

// ClaimedTask описывает задачу и конкретную попытку, которой runner выдал lease.
type ClaimedTask struct {
	// Исходные данные выполняемой задачи.
	Task `json:"task"`

	// Номер текущей попытки, используемый как fencing token.
	Attempt int `json:"attempt"`
}

type handler func(context.Context, json.RawMessage) error

// Descriptor связывает source задачи с типом её входных данных.
type Descriptor[TPayload any] struct {
	// Source определяет тип задачи и зарегистрированный для неё обработчик.
	Source string
}

// NewDescriptor создаёт типизированное описание задачи с указанным source.
func NewDescriptor[TPayload any](source string) Descriptor[TPayload] {
	return Descriptor[TPayload]{Source: source}
}

// NewTask создаёт задачу указанного типа с обязательным идентификатором и сериализует её типизированную нагрузку в JSON.
// Nil serialKey не ограничивает параллельное выполнение, а nil runAt означает запуск без отложенного старта.
func NewTask[TPayload any](descriptor Descriptor[TPayload], sourceID string, serialKey *string, runAt *time.Time, payload TPayload) (Task, error) {
	if sourceID == "" {
		return Task{}, errors.New("task source id is empty")
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return Task{}, fmt.Errorf("marshal task payload: %w", err)
	}
	taskRunAt := time.Now()
	if runAt != nil {
		taskRunAt = *runAt
	}

	return Task{
		Source:    descriptor.Source,
		SourceID:  sourceID,
		SerialKey: serialKey,
		Payload:   data,
		RunAt:     taskRunAt,
	}, nil
}

// HandleFunc выполняет бизнес-сценарий задачи с типизированными входными данными.
type HandleFunc[TPayload any] = func(context.Context, TPayload) error

type retryAtError struct {
	err     error
	retryAt time.Time
}

func (e *retryAtError) Error() string {
	return e.err.Error()
}

func (e *retryAtError) Unwrap() error {
	return e.err
}

// RetryAt возвращает ошибку, которая просит runner повторить задачу не раньше retryAt.
func RetryAt(err error, retryAt time.Time) error {
	if err == nil {
		return nil
	}

	return &retryAtError{err: err, retryAt: retryAt}
}
