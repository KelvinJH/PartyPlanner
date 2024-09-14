package main

import (
	"fmt"
	"partyplanner/db"
	router "partyplanner/router"
	service "partyplanner/service"
	ws "partyplanner/ws"
	"log"
	"net/http"
)

func main() {
	fmt.Println("This is a Calendar App on the web")
	setupHandlers()
	db.InitDatabase()
	database := db.GetDbInstance()
	defer database.CloseConnect()

	log.Fatal(http.ListenAndServe(":8008", nil))
}

func setupHandlers() {
	authPage := service.CreateAuthPage()

	manager := ws.NewManager()
	go manager.Run()

	// Pages 
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		authPage.Execute(w, nil)
	})
	http.HandleFunc("/calendar", func(w http.ResponseWriter, r *http.Request) {
		template, templateData := service.CreateCalendar()
		data := templateData
		template.Execute(w, data)
	})

	//Endpoints
	http.HandleFunc("/v1/event", router.SaveEvent)
	http.HandleFunc("/v1/authorize", router.AuthorizeUser)
	http.HandleFunc("/v1/room", router.CreateRoom)
	// Websocket
	http.HandleFunc("/wschat", manager.ServeWebsocket)
}
