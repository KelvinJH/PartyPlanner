package router

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"partyplanner/bus"
	"partyplanner/service"
	"strconv"

	"github.com/gorilla/sessions"
)

type Router struct {
	store    *sessions.CookieStore
	eventBus *bus.EventBus
}

type EventDto struct {
	Name        string `json:"event_name"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Description string `json:"event_description"`
}

type Secret struct {
	Key string `json:"secret_key"`
}

var router Router

func NewRouter(bus *bus.EventBus) *Router {
	secretJson, err := os.Open("./secret.json")
	if err != nil {
		fmt.Println("Failed to get secret key from json")
	}
	defer secretJson.Close()
	var session Secret
	secretBytes, _ := io.ReadAll(secretJson)
	json.Unmarshal(secretBytes, &session)

	router = Router{
		store:    sessions.NewCookieStore([]byte(session.Key)),
		eventBus: bus,
	}
	return &router
}

func SaveEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	event := EventDto{
		Name:        r.FormValue("event-name"),
		StartDate:   r.FormValue("start-date"),
		EndDate:     r.FormValue("end-date"),
		Description: r.FormValue("event-description"),
	}
	eventBytes, err := json.Marshal(event)
	if err != nil {
		fmt.Println("Error converting event to bytes: ", err)
		return
	}
	router.eventBus.Publish(eventBytes)
	session, _ := router.store.Get(r, "session.id")
	roomId, ok := session.Values["room-id"].(int)
	if !ok {
		http.Error(w, "Unable to find room id from session", http.StatusBadRequest)
		return
	}

	// Save to db after testing
	eventId, err := service.SaveEvent(roomId, event.Name, event.StartDate, event.EndDate, event.Description)
	if err != nil {
		http.Error(w, "Unable to save event", http.StatusInternalServerError)
		return
	}

	log.Println(eventId, " Has been saved")
}

func CreateRoom(w http.ResponseWriter, r *http.Request) {
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

	roomId := service.CreateRoom(r.FormValue("room-name"), r.FormValue("room-key"), int16(cap))
	if roomId != 0 {
		fmt.Printf("Created room %d successfully", roomId)
		session, _ := router.store.Get(r, "session.id")
		session.Values["authenticated"] = true
		session.Values["room-id"] = roomId
		session.Values["room-name"] = r.FormValue("room-name")
		http.Redirect(w, r, "/calendar", http.StatusSeeOther)
	}
}

func AuthorizeUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	roomName := r.FormValue("room-name")
	roomKey := r.FormValue("room-key")
	foundId, foundName, err := service.ValidateRoom(roomName, roomKey)
	if err != nil {
		http.Error(w, "Unable to validate (room, key)", http.StatusBadRequest)
	}
	if foundId != 0 && foundName != "" {
		fmt.Println("the room exists")
		session, _ := router.store.Get(r, "session.id")
		session.Values["authenticated"] = true
		session.Values["room-name"] = foundName
		session.Values["room-id"] = foundId
		session.Save(r, w)
	}

	http.Redirect(w, r, "/calendar", http.StatusSeeOther)

}

func Healthcheck(w http.ResponseWriter, r *http.Request) {
	authenticated := isAuthenticated(r)
	if authenticated {
		w.Write([]byte("Session is authorized"))
		return
	} else {
		http.Error(w, "Unauthorized", http.StatusForbidden)
	}
}

func isAuthenticated(r *http.Request) bool {
	session, _ := router.store.Get(r, "session.id")
	authenticated := session.Values["authenticated"]
	fmt.Println(session.Values["room-name"])
	if authenticated != nil && authenticated != false {
		return true
	}
	return false
}

func Calendar(w http.ResponseWriter, r *http.Request) {
	template, templateData := service.CreateCalendar()
	session, _ := router.store.Get(r, "session.id")
	calendarName, ok := session.Values["room-name"].(string)
	if ok {
		templateData.CalendarName = calendarName
	}
	template.Execute(w, templateData)

	roomId, ok := session.Values["room-id"].(int)
	if !ok {
		return
	}

	go LoadCalendar(roomId)
}

func Home(w http.ResponseWriter, r *http.Request) {
	authPage := service.CreateAuthPage()
	authPage.Execute(w, nil)
}

func Authorized(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := router.store.Get(r, "session.id")
		authenticated, ok := session.Values["authenticated"].(bool)
		if ok && authenticated {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Unauthorized", http.StatusForbidden)
		}
	})
}

func LoadCalendar(roomId int) {
	for {
		select {
		case <-router.eventBus.Ready:
			goto LoadEvents
		}
	}
LoadEvents:
	fmt.Println("Inside load events")
	savedEvents, err := service.LoadCalendar(roomId)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := range savedEvents {
		router.eventBus.Publish([]byte(savedEvents[i]))
	}
}
