package commands

import (
	"context"
	"encoding/json"
	"firebase.google.com/go/auth"
	"log"
	"net/http"
	"strings"
	"technicalevent.sfeir.lu/events"
	"technicalevent.sfeir.lu/firebase"
	"technicalevent.sfeir.lu/gcloud"
)

type ChangeDescriptionCommand struct {
	Description string `json:"description"`
}

func changeSubjectDescription(entityId string, changeDescriptionCommand *ChangeDescriptionCommand, token *auth.Token) (*events.ChangeDescriptionEvent, error) {

	changeDescriptionEvent := events.NewChangeDescriptionEvent(
		changeDescriptionCommand.Description,
		entityId,
		token.UID,
	)

	msg := gcloud.NewMessage(changeDescriptionEvent)

	publishError := gcloud.PublishEvent(context.Background(), *msg)
	if publishError != nil {
		log.Fatalf("Publish descrionChanged events: %v", publishError)
	}
	return &changeDescriptionEvent, nil
}

func ChangeDescription(w http.ResponseWriter, r *http.Request) {

	id := getEntityId(r)
	changeDescriptionCommand := getCommand(r, new(ChangeDescriptionCommand))
	bearer := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
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

	event, err := changeSubjectDescription(id, changeDescriptionCommand.(*ChangeDescriptionCommand), token)
	if err != nil {
		log.Fatalf("Change Description command: %v", err)
	}

	data, err := json.Marshal(event)
	if err != nil {
		log.Fatalf("Marshal events: %v", err)
	}

	if _, err := w.Write(data); err != nil {
		log.Fatalf("Write events to response: %v", err)
	}
}
