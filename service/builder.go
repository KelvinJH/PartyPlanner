package service

import (
	"bytes"
	"fmt"
	"html/template"
	"time"
)

var Days = [7]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
var Months = [12]string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}

type CalendarData struct {
	CalendarName string
	CurrentMonth string
	Months       [12]string
	Year         int
	DayNames     [7]string
	Days         []int
}

type ChatMessageData struct {
	ClientName string
	Content    string
	Timestamp  string
}

func CreateCalendar() (*template.Template, CalendarData) {
	calendar, err := template.ParseFiles("./service/html/calendar.html")
	today := time.Now()
	days := getDaysInMonth(today)

	if err != nil {
		fmt.Printf("Error while rendering calendar: %s \n", err)
	}
	calendarData := CalendarData{
		CurrentMonth: today.Month().String(),
		Months:       Months,
		Year:         today.Year(),
		DayNames:     Days,
		Days:         days,
	}
	return calendar, calendarData

}

func CreateAuthPage() *template.Template {
	waitingRoom, err := template.ParseFiles("./service/html/waitingroom.html")

	if err != nil {
		fmt.Println("Error while parsing the waiting room template file: ", err)
	}

	return waitingRoom
}

func CreateChatMessages(clientName, content, timestamp string) []byte {
	chatMessage, err := template.ParseFiles("./service/html/message.html")
	if err != nil {
		fmt.Printf("Error while rendering chat message: %s \n", err)
	}

	chatData := ChatMessageData{
		ClientName: clientName,
		Content:    content,
		Timestamp:  timestamp,
	}
	var renderedMessage bytes.Buffer
	err = chatMessage.Execute(&renderedMessage, chatData)
	if err != nil {
		fmt.Printf("Error while parsing chat message: %s \n", err)
	}

	return renderedMessage.Bytes()
}

func getDaysInMonth(now time.Time) []int {
	daysInMonth := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()

	daysBefore := int(time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Weekday())

	days := make([]int, daysInMonth+daysBefore)
	for i := range days {
		days[i] = i - daysBefore + 1
	}

	return days
}
