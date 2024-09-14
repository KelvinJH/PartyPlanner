package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "root"
	password = "secret"
	dbname   = "gotripdb"
)

type Room struct {
	Id          int
	Capacity    int
	Name        string
	CreatedDate time.Time
	UpdatedDate time.Time
	Key         string
}

type Database struct {
	driver *sql.DB
}

var database *Database

// Initialize the database instance and store it into the global var
func InitDatabase() {
	database = newDatabase()
}

// Get instance so that you can use it in other files
func GetDbInstance() *Database {
	return database
}

func newDatabase() *Database {
	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user, password, host, port, dbname)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatalf("Error opening DB connection: %v", err)
	}

	return &Database{
		driver: db,
	}
}

func (db *Database) Ping() {

	err := db.driver.Ping()

	if err != nil {
		log.Fatal("Error connecting to DB: ", err)
	}
}

func (db *Database) CloseConnect() {
	fmt.Println("Closing connection to db")
	db.driver.Close()
}

func (db *Database) GetCalendar(calendarId int) {
	fmt.Printf("Getting calendar %d from database\n", calendarId)
}

func (db *Database) CreateRoom(name, key string, cap int16) int {
	fmt.Printf("Creating %s with a capacity of %d \n", name, cap)


	// Check if room already exists 
	existingRoom := db.GetRoom(name, key)
	if existingRoom.Id != 0 {
		fmt.Printf("There is an existing room (%d - %s) with this name and key \n", existingRoom.Id, name)
		return existingRoom.Id
	}

	query := `INSERT INTO rooms (room_name, room_key, capacity)
		VALUES ($1, $2, $3) RETURNING room_id`

	var primaryKey int

	err := db.driver.QueryRow(query, name, key, cap).Scan(&primaryKey)
	if err != nil {
		log.Fatalf("Error while creating a new room %v\n", err)
	}

	return primaryKey
}

func (db *Database) GetRoom(name, key string) Room {
	query := `SELECT room_id, capacity, room_name FROM rooms WHERE room_name = $1 and room_key = $2`

	var foundRoom Room

	err := db.driver.QueryRow(query, name, key).Scan(&foundRoom.Id, &foundRoom.Capacity, &foundRoom.Name)
	if err != nil {
		log.Fatalf("Error while searching for room %s: %v\n", name, err)
	}
	return foundRoom
}
