package router

import (
	"fmt"
	service "partyplanner/service"
	"log"
	"net/http"
	"strconv"
)

type Event struct {
	Name        string `json:"event_name"`
	Date        string `json:"event_date"`
	Description string `json:"event_description"`
}

func SaveEvent(w http.ResponseWriter, r *http.Request) {
	log.Println("Received Request")

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	event := Event{
		Name:        r.FormValue("event-name"),
		Date:        r.FormValue("event-date"),
		Description: r.FormValue("event-description"),
	}

	log.Println(event, " Has been saved")
}

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	log.Println("Room created --> Redirecting you to calendar")

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	cap, err := strconv.Atoi(r.FormValue("room-capacity"))
	if err != nil {
		http.Error(w, "Unable to parse capacity from form", http.StatusBadRequest)
	}


	service.CreateRoom(r.FormValue("room-name"), r.FormValue("room-key"), int16(cap))
	http.Redirect(w, r, "/calendar", http.StatusSeeOther)
}

func AuthorizeUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Inside authorization workflow")

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	fmt.Printf("Room details : %s - %s and your name %s", r.FormValue("room-name"), r.FormValue("room-key"), r.FormValue("client-name"))

	roomName, err := service.ValidateRoom(r.FormValue("room-name"), r.FormValue("room-key"))
	if err != nil {
		http.Error(w, "Unable to validate (room, key)", http.StatusBadRequest)
	}
	fmt.Println("Found room: ", roomName)
	// if user is authorized
	http.Redirect(w, r, "/calendar", http.StatusSeeOther)

	// Consume token from user account ->
	// join room + load calendar page for them then load page
	//
}
