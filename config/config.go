package config

import "os"

var (
	Token    = ""
	DataFile = "tracked_users.json"
)


func LoadConfig() {
	Token = os.Getenv("DISCORD_TOKEN")
	if Token == "" {
		panic("❌ DISCORD_TOKEN non défini dans les variables d'environnement")
	}
}