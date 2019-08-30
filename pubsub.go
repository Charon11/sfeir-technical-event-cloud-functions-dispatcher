package subject

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"log"
)

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

func PublishMessage(ctx context.Context, m pubsub.Message, topic string) error {

	id, err := client.Topic(topic).Publish(ctx, &m).Get(ctx)
	if err != nil {
		log.Printf("topic(%s).Publish.Get: %v", topic, err)
		return err
	}
	fmt.Sprintf("Published msg: %s", id)
	log.Printf("Published msg: %s", id)
	return nil
}

func PublishEvent(ctx context.Context, m pubsub.Message) error {
	return PublishMessage(ctx, m, "subject-events")
}
