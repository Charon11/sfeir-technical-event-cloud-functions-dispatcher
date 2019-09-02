package events

type RefuseEvent struct {
	Event
	Status string `json:"status"`
}

func NewRefuseEvent(entityId string, userId string) RefuseEvent {
	return RefuseEvent{
		Status: "Refusé",
		Event:  New(entityId, "refused", userId),
	}
}
