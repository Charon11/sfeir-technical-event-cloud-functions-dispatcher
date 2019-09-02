package events

type ChangeDescriptionEvent struct {
	Description string `json:"description"`
	Event
}

func NewChangeDescriptionEvent(description string, entityId string, userId string) ChangeDescriptionEvent {
	return ChangeDescriptionEvent{
		Description: description,
		Event:       New(entityId, "descriptionChanged", userId),
	}
}
