package subject

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"firebase.google.com/go/auth"
	"github.com/google/uuid"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type DeleteEvent struct {
	Status   string    `json:"status"`
	Id       string    `json:"id"`
	Ts       time.Time `json:"_ts"`
	EntityId string    `json:"entityId"`
	UserId   string    `json:"userId"`
	Kind     string    `json:"_kind"`
}

func deleteEvent(entityId string, token *auth.Token) (*DeleteEvent, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		log.Fatalf("Generate uuid: %v", err)
		return nil, err
	}

	deleteEvent := DeleteEvent{
		Status:   "Supprimé",
		Id:       id.String(),
		Ts:       time.Now(),
		EntityId: entityId,
		UserId:   token.UID,
		Kind:     "deleted",
	}

	data, err := json.Marshal(deleteEvent)
	if err != nil {
		log.Fatalf("Marshal event: %v", err)
		return nil, err
	}
	msg := &pubsub.Message{
		Data:       data,
		Attributes: map[string]string{"_kind": deleteEvent.Kind},
	}

	publishError := PublishEvent(context.Background(), *msg)
	if publishError != nil {
		log.Fatalf("Publish deleted event: %v", publishError)
		return nil, publishError
	}
	return &deleteEvent, nil
}

func Delete(w http.ResponseWriter, r *http.Request) {

	re := regexp.MustCompile(`\w{8}-\w{4}-\w{4}-\w{4}-\w{12}`)
	id := re.FindString(r.URL.Path)
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

	event, err := deleteEvent(id, token)
	if err != nil {
		log.Fatalf("Delete command: %v", err)
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
