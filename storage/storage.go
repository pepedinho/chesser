package storage

import (
	"chesser/config"
	"encoding/json"
	"fmt"
	"os"
)

var TrackedUsers = make(map[string]string)

func LoadTrackedUser() {
	file, err := os.ReadFile(config.DataFile)
	if err == nil {
		json.Unmarshal(file, &TrackedUsers)
		fmt.Println("✅ Tracked user loaded succesfuly : ", TrackedUsers)
	}
}

func SaveTrackedUsers() {
	data, _ := json.MarshalIndent(TrackedUsers, "", "  ")
	os.WriteFile(config.DataFile, data, 0644)
}
