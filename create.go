package subject

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"firebase.google.com/go/auth"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"
	"time"
)

type SubjectType string

type CreatedEvent struct {
	Title       string      `json:"title"`
	Description string      `json:"description"`
	SubjectType SubjectType `json:"subjectType"`
	Record      bool        `json:"record"`
	Status      string      `json:"status"`
	Id          string      `json:"id"`
	Ts          time.Time   `json:"_ts"`
	EntityId    string      `json:"entityId"`
	UserId      string      `json:"userId"`
	Kind        string      `json:"_kind"`
}

type CreateCommand struct {
	Title       string      `json:"title"`
	Description string      `json:"description"`
	SubjectType SubjectType `json:"subjectType"`
	Record      bool        `json:"record"`
}

func createSubject(createCommand CreateCommand, token *auth.Token) (*CreatedEvent, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		log.Fatalf("Generate uuid: %v", err)
		return nil, err
	}

	entityId, err := uuid.NewRandom()
	if err != nil {
		log.Fatalf("Generate uuid: %v", err)
		return nil, err
	}

	createdEvent := CreatedEvent{
		Title:       createCommand.Title,
		Description: createCommand.Description,
		SubjectType: createCommand.SubjectType,
		Record:      createCommand.Record,
		Status:      "Nouveau",
		Id:          id.String(),
		Ts:          time.Now(),
		EntityId:    entityId.String(),
		UserId:      token.UID,
		Kind:        "created",
	}

	data, err := json.Marshal(createdEvent)
	if err != nil {
		log.Fatalf("Marshal event: %v", err)
		return nil, err
	}
	msg := &pubsub.Message{
		Data:       data,
		Attributes: map[string]string{"_kind": createdEvent.Kind},
	}

	publishError := PublishEvent(context.Background(), *msg)
	if publishError != nil {
		log.Fatalf("Publish created event: %v", publishError)
		return nil, publishError
	}
	return &createdEvent, nil
}

func Create(w http.ResponseWriter, r *http.Request) {

	var createCommand CreateCommand
	if err := json.NewDecoder(r.Body).Decode(&createCommand); err != nil {
		log.Fatalf("Unmarshal command: %v", err)
		return
	}
	bearer := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if len(bearer) < 1 {
		w.WriteHeader(401)
		return
	}
	token, err := VerifyIDToken(context.Background(), bearer)
	if err != nil {
		w.WriteHeader(403)
		_, _ = w.Write([]byte(err.Error()))
		log.Printf("Verify token: %v", err)
		return
	}

	event, err := createSubject(createCommand, token)
	if err != nil {
		log.Fatalf("Create command: %v", err)
	}

	data, err := json.Marshal(event)
	if err != nil {
		log.Fatalf("Marshal event: %v", err)
	}

	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		log.Fatalf("Write event to response: %v", err)
	}
}
