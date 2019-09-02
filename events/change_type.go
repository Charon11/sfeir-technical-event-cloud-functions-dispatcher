package events

type ChangeTypeEvent struct {
	SubjectType SubjectType `json:"subjectType"`
	Event
}

func NewChangeTypeEvent(subjectType SubjectType, entityId string, userId string) ChangeTypeEvent {
	return ChangeTypeEvent{
		SubjectType: subjectType,
		Event:       New(entityId, "typeChanged", userId),
	}
}
