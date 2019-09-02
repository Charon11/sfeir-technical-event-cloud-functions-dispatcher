package events

type CreatedEvent struct {
	Title       string      `json:"title"`
	Description string      `json:"description"`
	SubjectType SubjectType `json:"subjectType"`
	Record      bool        `json:"record"`
	Status      string      `json:"status"`
	Event
}

func NewCreatedEvent(
	title string,
	description string,
	subjectType SubjectType,
	record bool,
	entityId string,
	userId string,
) CreatedEvent {
	return CreatedEvent{
		Title:       title,
		Description: description,
		SubjectType: subjectType,
		Record:      record,
		Status:      "Nouveau",
		Event:       New(entityId, "created", userId),
	}
}
