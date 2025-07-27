package bot

import (
	"chesser/chess"
	"chesser/config"
	"chesser/roles"
	"chesser/storage"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var BotSession *discordgo.Session

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "chess",
		Description: "Associe ton compte Chess.com pour que le bot puisse tracker ton Elo",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "username",
				Description: "Ton pseudo Chess.com",
				Required:    true,
			},
		},
	},
}

func Start() error {
	var err error
	BotSession, err = discordgo.New("Bot " + config.Token)
	if err != nil {
		return err
	}

	BotSession.AddHandler(interactionHandler)

	err = BotSession.Open()
	if err != nil {
		return err
	}

	for _, cmd := range commands {
		_, err := BotSession.ApplicationCommandCreate(BotSession.State.User.ID, "", cmd)
		if err != nil {
			fmt.Println("❌ Erreur enregistrement commande :", err)
		}
	}

	go StartScheduler()

	fmt.Println("✅ Slash commands enregistrées avec succès")
	return nil
}

func interactionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	switch i.ApplicationCommandData().Name {
	case "chess":
		username := i.ApplicationCommandData().Options[0].StringValue()

		storage.TrackedUsers[i.Member.User.ID] = username
		storage.SaveTrackedUsers()

		// assignation du role chess si il ne l'a pas deja
		defaultRoleID := "1399078502334070874"
		member := i.Member
		hasDefaultRole := false
		for _, rID := range member.Roles {
			if rID == defaultRoleID {
				hasDefaultRole = true
				break
			}
		}
		if !hasDefaultRole {
			err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, defaultRoleID)
			if err != nil {
				fmt.Println("⚠️ Impossible d'ajouter le rôle par défaut :", err)
			}
		}

		rating, err := chess.FetchChessRating(username)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "❌ Erreur lors de la récupération du rating.",
				},
			})
			return
		}

		err = roles.UpdateUserRole(s, i.GuildID, i.Member.User.ID, rating)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "❌ Erreur lors de la mise à jour du rôle.",
				},
			})
			return
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("✅ Ton rôle a été mis à jour (%d Elo).", rating),
			},
		})
	}
}
