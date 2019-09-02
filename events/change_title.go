package events

type ChangeTitleEvent struct {
	Title string `json:"title"`
	Event
}

func NewChangeTitleEvent(title string, entityId string, userId string) ChangeTitleEvent {
	return ChangeTitleEvent{
		Title: title,
		Event: New(entityId, "titleChanged", userId),
	}
}
