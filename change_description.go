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

type ChangeDescriptionEvent struct {
	Description string    `json:"description"`
	Id          string    `json:"id"`
	Ts          time.Time `json:"_ts"`
	EntityId    string    `json:"entityId"`
	UserId      string    `json:"userId"`
	Kind        string    `json:"_kind"`
}

type ChangeDescriptionCommand struct {
	Description string `json:"description"`
}

func changeSubjectDescription(entityId string, changeDescriptionCommand ChangeDescriptionCommand, token *auth.Token) (*ChangeDescriptionEvent, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		log.Fatalf("Generate uuid: %v", err)
		return nil, err
	}
	changeDescriptionEvent := ChangeDescriptionEvent{
		Description: changeDescriptionCommand.Description,
		Id:          id.String(),
		Ts:          time.Now(),
		EntityId:    entityId,
		UserId:      token.UID,
		Kind:        "descriptionChanged",
	}

	data, err := json.Marshal(changeDescriptionEvent)
	if err != nil {
		log.Fatalf("Marshal event: %v", err)
		return nil, err
	}
	msg := &pubsub.Message{
		Data:       data,
		Attributes: map[string]string{"_kind": changeDescriptionEvent.Kind},
	}

	publishError := PublishEvent(context.Background(), *msg)
	if publishError != nil {
		log.Fatalf("Publish descrionChanged event: %v", publishError)
		return nil, publishError
	}
	return &changeDescriptionEvent, nil
}

func ChangeDescription(w http.ResponseWriter, r *http.Request) {

	re := regexp.MustCompile(`\w{8}-\w{4}-\w{4}-\w{4}-\w{12}`)
	id := re.FindString(r.URL.Path)
	var changeDescriptionCommand ChangeDescriptionCommand
	if err := json.NewDecoder(r.Body).Decode(&changeDescriptionCommand); err != nil {
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

	event, err := changeSubjectDescription(id, changeDescriptionCommand, token)
	if err != nil {
		log.Fatalf("Change Description command: %v", err)
	}

	data, err := json.Marshal(event)
	if err != nil {
		log.Fatalf("Marshal event: %v", err)
	}

	if _, err := w.Write(data); err != nil {
		log.Fatalf("Write event to response: %v", err)
	}
}
