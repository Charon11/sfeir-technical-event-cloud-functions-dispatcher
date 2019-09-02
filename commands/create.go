package commands

import (
	"context"
	"encoding/json"
	"firebase.google.com/go/auth"
	"github.com/google/uuid"
	"log"
	"net/http"
	"technicalevent.sfeir.lu/events"
	"technicalevent.sfeir.lu/firebase"
	"technicalevent.sfeir.lu/gcloud"
)

type CreateCommand struct {
	Title       string             `json:"title"`
	Description string             `json:"description"`
	SubjectType events.SubjectType `json:"subjectType"`
	Record      bool               `json:"record"`
}

func createSubject(createCommand *CreateCommand, token *auth.Token) (*events.CreatedEvent, error) {

	entityId, err := uuid.NewRandom()
	if err != nil {
		log.Fatalf("Generate uuid: %v", err)
	}

	createdEvent := events.NewCreatedEvent(
		createCommand.Title,
		createCommand.Description,
		createCommand.SubjectType,
		createCommand.Record,
		entityId.String(),
		token.UID,
	)

	msg := gcloud.NewMessage(createdEvent)

	publishError := gcloud.PublishEvent(context.Background(), *msg)
	if publishError != nil {
		log.Fatalf("Publish created events: %v", publishError)
	}
	return &createdEvent, nil
}

func Create(w http.ResponseWriter, r *http.Request) {

	createCommand := getCommand(r, new(CreateCommand))

	bearer := getBearerToken(r)
	if len(bearer) < 1 {
		w.WriteHeader(401)
		return
	}
	token, err := firebase.VerifyIDToken(context.Background(), bearer)
	if err != nil {
		w.WriteHeader(403)
		_, _ = w.Write([]byte(err.Error()))
		log.Printf("Verify token: %v", err)
		return
	}

	event, err := createSubject(createCommand.(*CreateCommand), token)
	if err != nil {
		log.Fatalf("Create command: %v", err)
	}

	data, err := json.Marshal(event)
	if err != nil {
		log.Fatalf("Marshal events: %v", err)
	}

	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		log.Fatalf("Write events to response: %v", err)
	}
}
