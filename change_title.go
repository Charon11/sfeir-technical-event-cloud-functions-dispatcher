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

type ChangeTitleEvent struct {
	Title    string    `json:"title"`
	Id       string    `json:"id"`
	Ts       time.Time `json:"_ts"`
	EntityId string    `json:"entityId"`
	UserId   string    `json:"userId"`
	Kind     string    `json:"_kind"`
}

type ChangeTitleCommand struct {
	Title string `json:"title"`
}

func changeSubjectTitle(entityId string, changeTitleCommand ChangeTitleCommand, token *auth.Token) (*ChangeTitleEvent, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		log.Fatalf("Generate uuid: %v", err)
		return nil, err
	}
	changeTitleEvent := ChangeTitleEvent{
		Title:    changeTitleCommand.Title,
		Id:       id.String(),
		Ts:       time.Now(),
		EntityId: entityId,
		UserId:   token.UID,
		Kind:     "titleChanged",
	}

	data, err := json.Marshal(changeTitleEvent)
	if err != nil {
		log.Fatalf("Marshal event: %v", err)
		return nil, err
	}
	msg := &pubsub.Message{
		Data:       data,
		Attributes: map[string]string{"_kind": changeTitleEvent.Kind},
	}

	publishError := PublishEvent(context.Background(), *msg)
	if publishError != nil {
		log.Fatalf("Publish titleChanged event: %v", publishError)
		return nil, publishError
	}
	return &changeTitleEvent, nil
}

func ChangeTitle(w http.ResponseWriter, r *http.Request) {

	re := regexp.MustCompile(`\w{8}-\w{4}-\w{4}-\w{4}-\w{12}`)
	id := re.FindString(r.URL.Path)
	var changeTitleCommand ChangeTitleCommand
	if err := json.NewDecoder(r.Body).Decode(&changeTitleCommand); err != nil {
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

	event, err := changeSubjectTitle(id, changeTitleCommand, token)
	if err != nil {
		log.Fatalf("change title command: %v", err)
	}

	data, err := json.Marshal(event)
	if err != nil {
		log.Fatalf("Marshal event: %v", err)
	}

	if _, err := w.Write(data); err != nil {
		log.Fatalf("Write event to response: %v", err)
	}
}
