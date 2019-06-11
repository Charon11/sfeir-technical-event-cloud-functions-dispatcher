// Package helloworld provides a set of Cloud Functions samples.
package subject

import (
	"context"
	"fmt"
	"log"
	"net/http"
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



// SubjectPubSub consumes a Pub/Sub message.
func SubjetPubSub(ctx context.Context, m pubsub.Message) error {
	id, err := client.Topic("subject-test").Publish(ctx, m).Get(r.Context())
	if err != nil {
		log.Printf("topic(%s).Publish.Get: %v", p.Topic, err)
		http.Error(w, "Error publishing message", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Published msg: %v", id)
}