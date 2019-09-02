package events

import "time"

type AcceptedEvent struct {
	Status       string    `json:"status"`
	Link         string    `json:"link"`
	ScheduleDate time.Time `json:"scheduleDate"`
	Event
}

func NewAcceptedEvent(scheduleDate time.Time, link string, entityId string, userId string) AcceptedEvent {
	return AcceptedEvent{
		Status:       "Accepté",
		Link:         link,
		ScheduleDate: scheduleDate,
		Event:        New(entityId, "accepted", userId),
	}
}
