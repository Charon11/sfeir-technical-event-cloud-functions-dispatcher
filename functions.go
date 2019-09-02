package functions

import (
	"cloud.google.com/go/pubsub"
	"context"
	"net/http"
	"technicalevent.sfeir.lu/commands"
	"technicalevent.sfeir.lu/gcloud"
)

// SubjectPubSub consumes a Pub/Sub message.
func SubjectPubSubDispatcher(ctx context.Context, m pubsub.Message) error {
	switch m.Attributes["_kind"] {
	case "accepted":
		return gcloud.PublishMessage(ctx, m, "accepted-subject-events")
	case "created":
		return gcloud.PublishMessage(ctx, m, "created-subject-events")
	case "refused":
		return gcloud.PublishMessage(ctx, m, "refused-subject-events")
	case "deleted":
		return gcloud.PublishMessage(ctx, m, "deleted-subject-events")
	case "descriptionChanged":
		return gcloud.PublishMessage(ctx, m, "description-changed-subject-events")
	case "schedulesChanged":
		return gcloud.PublishMessage(ctx, m, "schedules-changed-subject-events")
	case "titleChanged":
		return gcloud.PublishMessage(ctx, m, "title-changed-subject-events")
	case "typeChanged":
		return gcloud.PublishMessage(ctx, m, "type-changed-subject-events")
	case "recordAuthorisationChanged":
		return gcloud.PublishMessage(ctx, m, "record-auth-changed-subject-events")
	}
	return nil
}

func Accept(w http.ResponseWriter, r *http.Request) {
	commands.Accept(w, r)
}
func ChangeDescription(w http.ResponseWriter, r *http.Request) {
	commands.ChangeDescription(w, r)
}
func ChangeRecordAuthorisation(w http.ResponseWriter, r *http.Request) {
	commands.ChangeRecordAuthorisation(w, r)
}
func ChangeTitle(w http.ResponseWriter, r *http.Request) {
	commands.ChangeTitle(w, r)
}
func ChangeType(w http.ResponseWriter, r *http.Request) {
	commands.ChangeType(w, r)
}
func Create(w http.ResponseWriter, r *http.Request) {
	commands.Create(w, r)
}
func Delete(w http.ResponseWriter, r *http.Request) {
	commands.Delete(w, r)
}
func Refuse(w http.ResponseWriter, r *http.Request) {
	commands.Refuse(w, r)
}
