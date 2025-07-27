package bot

import (
	"chesser/chess"
	"chesser/roles"
	"chesser/storage"
	"fmt"
	"time"
)

func StartScheduler() {
	ticker := time.NewTicker(15 * time.Minute)
	for range ticker.C {
		if len(BotSession.State.Guilds) == 0 {
			continue
		}
		guildID := BotSession.State.Guilds[0].ID
		for discordID, chessUsername := range storage.TrackedUsers {
			rating, err := chess.FetchChessRating(chessUsername)
			if err != nil {
				fmt.Println("❌ API error for {", chessUsername, "} : ", err)
				continue
			}
			roles.UpdateUserRole(BotSession, guildID, discordID, rating)
		}
	}
}
