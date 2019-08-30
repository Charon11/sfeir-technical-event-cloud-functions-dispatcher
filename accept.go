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

type AcceptedEvent struct {
	Status       string    `json:"status"`
	Link         string    `json:"link"`
	ScheduleDate time.Time `json:"scheduleDate"`
	Id           string    `json:"id"`
	Ts           time.Time `json:"_ts"`
	EntityId     string    `json:"entityId"`
	UserId       string    `json:"userId"`
	Kind         string    `json:"_kind"`
}

type AcceptCommand struct {
	ScheduleDate time.Time `json:"scheduleDate"`
	Link         string    `json:"link"`
}

func acceptSubject(entityId string, acceptCommand AcceptCommand, token *auth.Token) (*AcceptedEvent, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		log.Fatalf("Generate uuid: %v", err)
		return nil, err
	}
	acceptedEvent := AcceptedEvent{
		Status:       "Accepté",
		ScheduleDate: acceptCommand.ScheduleDate,
		Link:         acceptCommand.Link,
		Id:           id.String(),
		Ts:           time.Now(),
		EntityId:     entityId,
		UserId:       token.UID,
		Kind:         "accepted",
	}

	data, err := json.Marshal(acceptedEvent)
	if err != nil {
		log.Fatalf("Marshal event: %v", err)
		return nil, err
	}
	msg := &pubsub.Message{
		Data:       data,
		Attributes: map[string]string{"_kind": acceptedEvent.Kind},
	}

	publishError := PublishEvent(context.Background(), *msg)
	if publishError != nil {
		log.Fatalf("Publish accepted event: %v", publishError)
		return nil, publishError
	}
	return &acceptedEvent, nil
}

func Accept(w http.ResponseWriter, r *http.Request) {

	re := regexp.MustCompile(`\w{8}-\w{4}-\w{4}-\w{4}-\w{12}`)
	id := re.FindString(r.URL.Path)

	var acceptCommand AcceptCommand
	if err := json.NewDecoder(r.Body).Decode(&acceptCommand); err != nil {
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

	event, err := acceptSubject(id, acceptCommand, token)
	if err != nil {
		log.Fatalf("Accept command: %v", err)
	}

	data, err := json.Marshal(event)
	if err != nil {
		log.Fatalf("Marshal event: %v", err)
	}

	if _, err := w.Write(data); err != nil {
		log.Fatalf("Write event to response: %v", err)
	}
}
