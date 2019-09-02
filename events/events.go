package events

import (
	"github.com/google/uuid"
	"log"
	"time"
)

type SubjectType string

type EventType interface {
	GetKind() string
}

func (event Event) GetKind() string {
	return event.Kind
}

type Event struct {
	Id       string    `json:"id"`
	Ts       time.Time `json:"_ts"`
	EntityId string    `json:"entityId"`
	UserId   string    `json:"userId"`
	Kind     string    `json:"_kind"`
}

func New(entityId string, kind string, userId string) Event {
	id, err := uuid.NewRandom()
	if err != nil {
		log.Fatalf("Generate uuid: %v", err)
	}
	return Event{
		Id:       id.String(),
		Ts:       time.Now(),
		EntityId: entityId,
		UserId:   userId,
		Kind:     kind,
	}
}
