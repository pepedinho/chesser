package bot

import (
	"chesser/chess"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func BuildChessEmbed(username string, stats chess.ChessStats) (*discordgo.MessageEmbed, error) {
	userData, err := chess.FetchChessUser(username)

	if err != nil {
		return nil, err
	}

	embed := &discordgo.MessageEmbed{
		Title:       username,
		URL:         fmt.Sprintf("https://www.chess.com/member/%s", username),
		Color:       0x1E90FF,
		Description: "**Ton role a été mise à jours♟️**",

		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: userData.Avatar,
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Blitz",
				Value:  fmt.Sprintf("`%d`", stats.ChessBlitz.Last.Rating),
				Inline: true,
			},
			{
				Name:   "Bullet",
				Value:  fmt.Sprintf("`%d`", stats.ChessBullet.Last.Rating),
				Inline: true,
			},
			{
				Name:   "Rapid",
				Value:  fmt.Sprintf("`%d`", stats.ChessRapid.Last.Rating),
				Inline: true,
			},
			{
				Name:   "Peak elo",
				Value:  fmt.Sprintf("**`%d`**", stats.ChessRapid.Best.Rating),
				Inline: false,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	return embed, nil
}
