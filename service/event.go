package service

import (
	"encoding/json"
	"fmt"
	"partyplanner/db"
	"time"
)

func SaveEvent(roomId int, name, start, end, description string) (int, error) {
	fmt.Println("Inside save")
	dbInstance := db.GetDbInstance()
	startDate, _ := time.Parse(db.TimeFormat, start)
	endDate, _ := time.Parse(db.TimeFormat, end)
	event := db.Event{
		RoomId:      roomId,
		Name:        name,
		StartDate:   startDate,
		EndDate:     endDate,
		Description: description,
	}
	fmt.Println(event)
	eventId, err := dbInstance.SaveEvent(event)
	if err != nil {
		return eventId, err
	}

	return eventId, nil
}

func LoadCalendar(roomId int) ([]string, error) {
	dbInstance := db.GetDbInstance()
	events, err := dbInstance.GetCalendar(roomId)
	if err != nil {
		fmt.Println(err)
		return make([]string, 0), err
	}
	savedEvents := make([]string, len(events))
	for i := range events {
		eventBytes, err := json.Marshal(events[i])
		if err != nil {
			fmt.Println("Unable to marshal events to bytes")
			return savedEvents, err
		}
		savedEvents = append(savedEvents, string(eventBytes))
	}
	return savedEvents, nil
}
