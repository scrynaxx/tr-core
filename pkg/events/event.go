package events

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Descriptor описывает тип события и содержит его имя для маршрутизации.
type Descriptor[TData any] struct {
	// Name — уникальное имя события, используемое как routing key.
	Name string
}

// NewDescriptor создаёт типизированный дескриптор события с заданным именем.
func NewDescriptor[TData any](name string) Descriptor[TData] {
	return Descriptor[TData]{Name: name}
}

// Message представляет событие с уже сериализованными данными.
type Message struct {
	// EventID однозначно идентифицирует конкретный экземпляр события.
	EventID uuid.UUID `json:"event_id"`

	// Timestamp содержит время возникновения события в UTC.
	Timestamp time.Time `json:"timestamp"`

	// Name определяет тип события и используется как routing key.
	Name string `json:"name"`

	// Data содержит сериализованную полезную нагрузку события.
	Data json.RawMessage `json:"data"`
}

// Event представляет типизированное событие, передаваемое обработчику.
type Event[TData any] struct {
	// EventID однозначно идентифицирует конкретный экземпляр события.
	EventID uuid.UUID `json:"event_id"`

	// Timestamp содержит время возникновения события в UTC.
	Timestamp time.Time `json:"timestamp"`

	// Name определяет тип события.
	Name string `json:"name"`

	// Data содержит типизированную полезную нагрузку события.
	Data TData `json:"data"`
}

// NewMessage создаёт событие, сериализуя данные согласно переданному дескриптору.
func NewMessage[TData any](desc Descriptor[TData], data TData) (Message, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return Message{}, fmt.Errorf("marshal event data: %w", err)
	}

	return Message{
		EventID:   uuid.New(),
		Timestamp: time.Now().UTC(),
		Name:      desc.Name,
		Data:      payload,
	}, nil
}
