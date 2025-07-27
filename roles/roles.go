package roles

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func getColorForElo(rating int) int {
	switch {
	case rating <= 450:
		return 0x9683EC
	case rating <= 600:
		return 0x361FDB
	case rating <= 800:
		return 0xF54927
	default:
		return 0xEFBF04
	}
}

func moveRoleAboveSpecific(s *discordgo.Session, guildID, roleID, targetRoleID string) error {
	roles, err := s.GuildRoles(guildID)
	if err != nil {
		return err
	}

	var targetRole, roleToMove *discordgo.Role
	for _, r := range roles {
		if r.ID == targetRoleID {
			targetRole = r
		}
		if r.ID == roleID {
			roleToMove = r
		}
	}

	if targetRole == nil {
		return fmt.Errorf("❌ Rôle cible (%s) introuvable", targetRoleID)
	}
	if roleToMove == nil {
		return fmt.Errorf("❌ Rôle à déplacer (%s) introuvable", roleID)
	}

	roleToMove.Position = targetRole.Position + 1

	_, err = s.GuildRoleReorder(guildID, roles)
	return err
}

func UpdateUserRole(s *discordgo.Session, guildID, userID string, newRating int) error {
	newRoleName := fmt.Sprintf("Chess (%d)", newRating)
	newRoleColor := getColorForElo(newRating)

	roles, err := s.GuildRoles(guildID)
	if err != nil {
		return err
	}

	var newRoleID string

	

	member, _ := s.GuildMember(guildID, userID)

	// si il y a aucun changement on ne fait rien
	for _, rID := range member.Roles {
		for _, role := range roles {
			if role.ID == rID && role.Name == newRoleName {
				//
				return nil
			}
		}
	}

	// suppresion du role conditionel
	for _, rID := range member.Roles {
		for _, role := range roles {
			if role.ID == rID && strings.HasPrefix(role.Name, "Chess (") {
				// Retirer l'ancien rôle
				s.GuildMemberRoleRemove(guildID, userID, rID)

				// Vérifier si quelqu'un d'autre utilise encore ce rôle
				stillUsed := false
				members, _ := s.GuildMembers(guildID, "", 1000)
				for _, m := range members {
					for _, mr := range m.Roles {
						if mr == rID {
							stillUsed = true
							break
						}
					}
					if stillUsed {
						break
					}
				}

				// Si plus personne ne l'utilise → on le supprime
				if !stillUsed {
					fmt.Println("🗑 Unused role has deleted succesfully :", role.Name)
					s.GuildRoleDelete(guildID, rID)
				}
			}
		}
	}

	for _, r := range roles {
		if r.Name == newRoleName {
			newRoleID = r.ID
			break
		}
	}

	// si le role existe pas deja on le creer
	if newRoleID == "" {
		hoist := false
		mentionable := true
		role, err := s.GuildRoleCreate(guildID, &discordgo.RoleParams{
			Name:        newRoleName,
			Color:       &newRoleColor,
			Hoist:       &hoist,
			Mentionable: &mentionable,
		})
		if err != nil {
			return err
		}
		if err = moveRoleAboveSpecific(s, guildID, role.ID, "1365115868899704962"); err != nil {
			fmt.Println("⚠️ Impossible de déplacer le rôle :", err)
		}
		newRoleID = role.ID
	}
	return s.GuildMemberRoleAdd(guildID, userID, newRoleID)
}
