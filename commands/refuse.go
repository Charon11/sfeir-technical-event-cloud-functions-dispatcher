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

func refuseEvent(entityId string, token *auth.Token) (*events.RefuseEvent, error) {

	refuseEvent := events.NewRefuseEvent(entityId, token.UID)

	msg := gcloud.NewMessage(refuseEvent)

	publishError := gcloud.PublishEvent(context.Background(), *msg)
	if publishError != nil {
		log.Fatalf("Publish refused events: %v", publishError)
	}
	return &refuseEvent, nil
}

func Refuse(w http.ResponseWriter, r *http.Request) {

	id := getEntityId(r)

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

	event, err := refuseEvent(id, token)
	if err != nil {
		log.Fatalf("Refuse command: %v", err)
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
