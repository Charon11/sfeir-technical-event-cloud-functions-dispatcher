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

type ChangeTitleCommand struct {
	Title string `json:"title"`
}

func changeSubjectTitle(entityId string, changeTitleCommand *ChangeTitleCommand, token *auth.Token) (*events.ChangeTitleEvent, error) {

	changeTitleEvent := events.NewChangeTitleEvent(
		changeTitleCommand.Title,
		entityId,
		token.UID,
	)

	msg := gcloud.NewMessage(changeTitleEvent)

	publishError := gcloud.PublishEvent(context.Background(), *msg)
	if publishError != nil {
		log.Fatalf("Publish titleChanged events: %v", publishError)
	}
	return &changeTitleEvent, nil
}

func ChangeTitle(w http.ResponseWriter, r *http.Request) {

	id := getEntityId(r)
	changeTitleCommand := getCommand(r, new(ChangeTitleCommand))
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

	event, err := changeSubjectTitle(id, changeTitleCommand.(*ChangeTitleCommand), token)
	if err != nil {
		log.Fatalf("change title command: %v", err)
	}

	data, err := json.Marshal(event)
	if err != nil {
		log.Fatalf("Marshal events: %v", err)
	}

	if _, err := w.Write(data); err != nil {
		log.Fatalf("Write events to response: %v", err)
	}
}
