package service

import (
	sha "crypto/sha256"
	"fmt"
	"partyplanner/db"
)

func ValidateRoom(name, key string) (bool, error) {
	fmt.Println("Inside the validate room key endpoint ---- Returning false for now")
	dbInstance := db.GetDbInstance()

	room := dbInstance.GetRoom(name, string(getHashedKey(key)))
	if room.Id == 0 {
		return false, fmt.Errorf("validation error")
	}
	return true, nil
}

func CreateRoom(name, key string, capacity int16) int {
	dbInstance := db.GetDbInstance()

	hashedKey := getHashedKey(key)
	createdRoom := dbInstance.CreateRoom(name, string(hashedKey), capacity)
	return createdRoom
}

func getHashedKey(key string) []byte {
	hash := sha.New()
	hash.Write([]byte(key))
	return hash.Sum(nil)
}
