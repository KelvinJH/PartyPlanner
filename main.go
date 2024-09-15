package main

import (
	"fmt"
	"log"
	"net/http"
	"partyplanner/bus"
	"partyplanner/db"
	"partyplanner/router"
	"partyplanner/ws"
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
	eventBus := bus.NewEventBus()
	chatManager := ws.NewMessageManager()
	eventManager := ws.NewEventManager(eventBus)
	router.NewRouter(eventBus)

	go chatManager.Run()
	go eventManager.ListenForEvents()
	// Pages
	http.HandleFunc("/", router.Authorized(router.Home))
	http.HandleFunc("/calendar", router.Authorized(router.Calendar))

	// Endpoints
	http.HandleFunc("/v1/event", router.Authorized(router.SaveEvent))
	http.HandleFunc("/v1/authorize", router.AuthorizeUser)
	http.HandleFunc("/v1/room", router.CreateRoom)
	http.HandleFunc("/v1/healthcheck", router.Healthcheck)

	// Websocket
	http.HandleFunc("/wschat", router.Authorized(chatManager.ServeChat))
	http.HandleFunc("/wsevent", router.Authorized(eventManager.ServeEvents))
}
