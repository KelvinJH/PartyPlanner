package service

import (
	"fmt"
	"partyplanner/db"
	"time"
)



func SaveEvent(roomId int, name, start, end, description string) (int, error) {
	dbInstance := db.GetDbInstance()
	fmt.Printf("start : %s \nend: %s\n", start, end)
	startDate, _ := time.Parse(db.TimeFormat, start)
	endDate, _ := time.Parse(db.TimeFormat, end)
	event := db.Event{
		RoomId: roomId,
		Name: name,
		StartDate: startDate,
		EndDate: endDate,
		Description: description,
	}
	fmt.Println(event)
	eventId, err := dbInstance.SaveEvent(event)
	if err != nil {
		return eventId, err
	}

	return eventId, nil
}
