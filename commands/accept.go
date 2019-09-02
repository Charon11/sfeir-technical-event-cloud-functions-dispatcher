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
	"time"
)

type AcceptCommand struct {
	ScheduleDate time.Time `json:"scheduleDate"`
	Link         string    `json:"link"`
}

func acceptSubject(entityId string, acceptCommand *AcceptCommand, token *auth.Token) (*events.AcceptedEvent, error) {
	acceptedEvent := events.NewAcceptedEvent(
		acceptCommand.ScheduleDate,
		acceptCommand.Link,
		entityId,
		token.UID,
	)

	msg := gcloud.NewMessage(acceptedEvent)

	publishError := gcloud.PublishEvent(context.Background(), *msg)
	if publishError != nil {
		log.Fatalf("Publish accepted events: %v", publishError)
	}
	return &acceptedEvent, nil
}

func Accept(w http.ResponseWriter, r *http.Request) {

	id := getEntityId(r)

	acceptCommand := getCommand(r, new(AcceptCommand))
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

	event, err := acceptSubject(id, acceptCommand.(*AcceptCommand), token)
	if err != nil {
		log.Fatalf("Accept command: %v", err)
	}

	data, err := json.Marshal(event)
	if err != nil {
		log.Fatalf("Marshal events: %v", err)
	}

	if _, err := w.Write(data); err != nil {
		log.Fatalf("Write events to response: %v", err)
	}
}
