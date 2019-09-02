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

type ChangeRecordAuthorisationCommand struct {
	Record bool `json:"record"`
}

func changeSubjectRecordAuthorisation(entityId string, changeRecordAuthorisationCommand *ChangeRecordAuthorisationCommand, token *auth.Token) (*events.ChangeRecordAuthorisationEvent, error) {
	changeRecordAuthorisationEvent := events.NewChangeRecordAuthorisationEvent(
		changeRecordAuthorisationCommand.Record,
		entityId,
		token.UID,
	)

	msg := gcloud.NewMessage(changeRecordAuthorisationEvent)

	publishError := gcloud.PublishEvent(context.Background(), *msg)
	if publishError != nil {
		log.Fatalf("Publish recordAuth change events: %v", publishError)
	}
	return &changeRecordAuthorisationEvent, nil
}

func ChangeRecordAuthorisation(w http.ResponseWriter, r *http.Request) {

	id := getEntityId(r)
	changeRecordAuthorisationCommand := getCommand(r, new(ChangeRecordAuthorisationCommand))
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

	event, err := changeSubjectRecordAuthorisation(id, changeRecordAuthorisationCommand.(*ChangeRecordAuthorisationCommand), token)
	if err != nil {
		log.Fatalf("Change RecordAuthorisation command: %v", err)
	}

	data, err := json.Marshal(event)
	if err != nil {
		log.Fatalf("Marshal events: %v", err)
	}

	if _, err := w.Write(data); err != nil {
		log.Fatalf("Write events to response: %v", err)
	}
}
