package events

type ChangeRecordAuthorisationEvent struct {
	Record bool `json:"record"`
	Event
}

func NewChangeRecordAuthorisationEvent(record bool, entityId string, userId string) ChangeRecordAuthorisationEvent {
	return ChangeRecordAuthorisationEvent{
		Record: record,
		Event:  New(entityId, "recordAuthorisationChanged", userId),
	}
}
