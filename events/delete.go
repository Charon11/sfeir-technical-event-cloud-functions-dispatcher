package events

type DeleteEvent struct {
	Status string `json:"status"`
	Event
}

func NewDeleteEvent(entityId string, userId string) DeleteEvent {
	return DeleteEvent{
		Status: "Supprimé",
		Event:  New(entityId, "deleted", userId),
	}
}
