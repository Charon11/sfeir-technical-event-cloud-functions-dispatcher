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

type ChangeTypeEvent struct {
	SubjectType SubjectType `json:"subjectType"`
	Id          string      `json:"id"`
	Ts          time.Time   `json:"_ts"`
	EntityId    string      `json:"entityId"`
	UserId      string      `json:"userId"`
	Kind        string      `json:"_kind"`
}

type ChangeTypeCommand struct {
	SubjectType SubjectType `json:"subjectType"`
}

func changeSubjectType(entityId string, changeTypeCommand ChangeTypeCommand, token *auth.Token) (*ChangeTypeEvent, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		log.Fatalf("Generate uuid: %v", err)
		return nil, err
	}
	changeTypeEvent := ChangeTypeEvent{
		SubjectType: changeTypeCommand.SubjectType,
		Id:          id.String(),
		Ts:          time.Now(),
		EntityId:    entityId,
		UserId:      token.UID,
		Kind:        "typeChanged",
	}

	data, err := json.Marshal(changeTypeEvent)
	if err != nil {
		log.Fatalf("Marshal event: %v", err)
		return nil, err
	}
	msg := &pubsub.Message{
		Data:       data,
		Attributes: map[string]string{"_kind": changeTypeEvent.Kind},
	}

	publishError := PublishEvent(context.Background(), *msg)
	if publishError != nil {
		log.Fatalf("Publish typeChanged event: %v", publishError)
		return nil, publishError
	}
	return &changeTypeEvent, nil
}

func ChangeType(w http.ResponseWriter, r *http.Request) {

	re := regexp.MustCompile(`\w{8}-\w{4}-\w{4}-\w{4}-\w{12}`)
	id := re.FindString(r.URL.Path)
	var changeTypeCommand ChangeTypeCommand
	if err := json.NewDecoder(r.Body).Decode(&changeTypeCommand); err != nil {
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

	event, err := changeSubjectType(id, changeTypeCommand, token)
	if err != nil {
		log.Fatalf("Change Type command: %v", err)
	}

	data, err := json.Marshal(event)
	if err != nil {
		log.Fatalf("Marshal event: %v", err)
	}

	if _, err := w.Write(data); err != nil {
		log.Fatalf("Write event to response: %v", err)
	}
}
