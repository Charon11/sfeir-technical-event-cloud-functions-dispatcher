package commands

import (
	"context"
	"encoding/json"
	"firebase.google.com/go/auth"
	"log"
	"net/http"
	"technicalevent.sfeir.lu/events"
	"technicalevent.sfeir.lu/firebase"
	"technicalevent.sfeir.lu/gcloud"
)

type ChangeTypeCommand struct {
	SubjectType events.SubjectType `json:"subjectType"`
}

func changeSubjectType(entityId string, changeTypeCommand *ChangeTypeCommand, token *auth.Token) (*events.ChangeTypeEvent, error) {

	changeTypeEvent := events.NewChangeTypeEvent(
		changeTypeCommand.SubjectType,
		entityId,
		token.UID,
	)

	msg := gcloud.NewMessage(changeTypeEvent)

	publishError := gcloud.PublishEvent(context.Background(), *msg)
	if publishError != nil {
		log.Fatalf("Publish typeChanged events: %v", publishError)
		return nil, publishError
	}
	return &changeTypeEvent, nil
}

func ChangeType(w http.ResponseWriter, r *http.Request) {

	id := getEntityId(r)
	changeTypeCommand := getCommand(r, new(ChangeTypeCommand))
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

	event, err := changeSubjectType(id, changeTypeCommand.(*ChangeTypeCommand), token)
	if err != nil {
		log.Fatalf("Change Type command: %v", err)
	}

	data, err := json.Marshal(event)
	if err != nil {
		log.Fatalf("Marshal events: %v", err)
	}

	if _, err := w.Write(data); err != nil {
		log.Fatalf("Write events to response: %v", err)
	}
}
