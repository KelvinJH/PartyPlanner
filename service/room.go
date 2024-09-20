package service

import (
	sha "crypto/sha256"
	"partyplanner/db"
)

func ValidateRoom(name, key string) (int, string, error) {
	dbInstance := db.GetDbInstance()

	room, err := dbInstance.GetRoom(name, string(getHashedKey(key)))
	if room.Id == 0 {
		return 0, "", err
	}
	return room.Id, room.Name, nil
}

func CreateRoom(name, key string, capacity int16) int {
	dbInstance := db.GetDbInstance()

	hashedKey := getHashedKey(key)
	createdRoom, err := dbInstance.CreateRoom(name, string(hashedKey), capacity)
	if err != nil {
		return 0
	}
	return createdRoom
}

func getHashedKey(key string) []byte {
	hash := sha.New()
	hash.Write([]byte(key))
	return hash.Sum(nil)
}
