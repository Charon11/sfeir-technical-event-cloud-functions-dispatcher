package gcloud

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"log"
	"os"
	"technicalevent.sfeir.lu/events"
)

var projectID = os.Getenv("GCP_PROJECT")

// client is a global Pub/Sub client, initialized once per instance.
var client *pubsub.Client

func init() {
	// err is pre-declared to avoid shadowing client.
	var err error
	// client is initialized with context.Background() because it should
	// persist between function invocations.
	client, err = pubsub.NewClient(context.Background(), projectID)
	if err != nil {
		log.Fatalf("pubsub.NewClient: %v", err)
	}
}

func NewMessage(event events.EventType) *pubsub.Message {
	data, err := json.Marshal(event)
	if err != nil {
		log.Fatalf("Marshal events: %v", err)
	}
	msg := &pubsub.Message{
		Data:       data,
		Attributes: map[string]string{"_kind": event.GetKind()},
	}
	return msg
}

func PublishMessage(ctx context.Context, m pubsub.Message, topic string) error {

	id, err := client.Topic(topic).Publish(ctx, &m).Get(ctx)
	if err != nil {
		log.Printf("topic(%s).Publish.Get: %v", topic, err)
		return err
	}
	log.Printf("Published msg: %s", id)
	return nil
}

func PublishEvent(ctx context.Context, m pubsub.Message) error {
	return PublishMessage(ctx, m, "subject-events")
}
