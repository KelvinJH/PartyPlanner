package router

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"io"
	"log"
	"net/http"
	"os"
	"partyplanner/service"
	"strconv"
)

var store *sessions.CookieStore
var session Secret

func InitRouter() {
	secretJson, err := os.Open("./secret.json")
	if err != nil {
		fmt.Println("Failed to get secret key from json")
	}
	defer secretJson.Close()

	secretBytes, _ := io.ReadAll(secretJson)
	json.Unmarshal(secretBytes, &session)
	store = sessions.NewCookieStore([]byte(session.Key))
}

type Event struct {
	Name        string `json:"event_name"`
	Date        string `json:"event_date"`
	Description string `json:"event_description"`
}

type Secret struct {
	Key string `json:"secret_key"`
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

	event := Event{
		Name:        r.FormValue("event-name"),
		Date:        r.FormValue("event-date"),
		Description: r.FormValue("event-description"),
	}

	log.Println(event, " Has been saved")
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
		session, _ := store.Get(r, "session.id")
		session.Values["room-name"] = r.FormValue("room-name")
		http.Redirect(w, r, "/calendar", http.StatusSeeOther)
	}
	return
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

	roomName := r.FormValue("room-name")
	roomKey := r.FormValue("room-key")
	exists, err := service.ValidateRoom(roomName, roomKey)
	if err != nil {
		http.Error(w, "Unable to validate (room, key)", http.StatusBadRequest)
	}
	if exists {
		fmt.Println("the room exists")
		session, _ := store.Get(r, "session.id")
		session.Values["authenticated"] = true
		session.Values["room-name"] = roomName
		session.Save(r, w)
	}

	http.Redirect(w, r, "/calendar", http.StatusSeeOther)

}

func Healthcheck(w http.ResponseWriter, r *http.Request) {
	authenticated := isAuthenticated(r)
	if authenticated == true {
		w.Write([]byte("Session is authorized"))
		return
	} else {
		http.Error(w, "Unauthorized", http.StatusForbidden)
	}
}

func isAuthenticated(r *http.Request) bool {
	session, _ := store.Get(r, "session.id")
	authenticated := session.Values["authenticated"]
	fmt.Println(session.Values["room-name"])
	if authenticated != nil && authenticated != false {
		return true
	}
	return false
}

func Calendar(w http.ResponseWriter, r *http.Request) {
	template, templateData := service.CreateCalendar()
	session, _ := store.Get(r, "session.id")
	calendarName, ok := session.Values["room-name"].(string)
	if ok {
		templateData.CalendarName = calendarName
	} else {
		templateData.CalendarName = "Unnamed"
	}
	template.Execute(w, templateData)
}

func Home(w http.ResponseWriter, r *http.Request) {
	authPage := service.CreateAuthPage()
	authPage.Execute(w, nil)
}

func Authorized(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session.id")
		authenticated := session.Values["authenticated"]
		if authenticated != nil && authenticated != false {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Unauthorized", http.StatusForbidden)
		}
	})
}
