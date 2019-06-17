package subject

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
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

func publishMessage(ctx context.Context, m pubsub.Message, topic string) error {
	id, err := client.Topic(topic).Publish(ctx, &m).Get(ctx)
	if err != nil {
		log.Printf("topic(%s).Publish.Get: %v", topic, err)
		return err
	}
	fmt.Sprintf("Published msg: %s", id)
	log.Printf("Published msg: %s", id)
	return nil
}


// SubjectPubSub consumes a Pub/Sub message.
func SubjectPubSubDispatcher(ctx context.Context, m pubsub.Message) error {
	switch m.Attributes["_kind"] {
	case "accepted" :
		return publishMessage(ctx, m, "accepted-subject-events")
	case "created" :
		return publishMessage(ctx, m, "created-subject-events")
	case "refused" :
		return publishMessage(ctx, m, "refused-subject-events")
	case "deleted" :
		return publishMessage(ctx, m, "deleted-subject-events")
	case "descriptionChanged" :
		return publishMessage(ctx, m, "description-changed-subject-events")
	case "schedulesChanged" :
		return publishMessage(ctx, m, "schedules-changed-subject-events")
	case "titleChanged" :
		return publishMessage(ctx, m, "title-changed-subject-events")
	}
	return nil
}