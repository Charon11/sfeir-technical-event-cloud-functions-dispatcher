package subject

import (
	"context"
	"os"

	"cloud.google.com/go/pubsub"
)

var projectID = os.Getenv("GCP_PROJECT")

// client is a global Pub/Sub client, initialized once per instance.
var client *pubsub.Client

// SubjectPubSub consumes a Pub/Sub message.
func SubjectPubSubDispatcher(ctx context.Context, m pubsub.Message) error {
	switch m.Attributes["_kind"] {
	case "accepted":
		return PublishMessage(ctx, m, "accepted-subject-events")
	case "created":
		return PublishMessage(ctx, m, "created-subject-events")
	case "refused":
		return PublishMessage(ctx, m, "refused-subject-events")
	case "deleted":
		return PublishMessage(ctx, m, "deleted-subject-events")
	case "descriptionChanged":
		return PublishMessage(ctx, m, "description-changed-subject-events")
	case "schedulesChanged":
		return PublishMessage(ctx, m, "schedules-changed-subject-events")
	case "titleChanged":
		return PublishMessage(ctx, m, "title-changed-subject-events")
	case "typeChanged":
		return PublishMessage(ctx, m, "type-changed-subject-events")
	case "recordAuthorisationChanged":
		return PublishMessage(ctx, m, "record-auth-changed-subject-events")
	}
	return nil
}
