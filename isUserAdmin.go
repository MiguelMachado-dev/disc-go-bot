package main

import "github.com/bwmarrin/discordgo"

func isUserAdmin(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	guildRoles, _ := s.GuildRoles(m.GuildID)
	member, _ := s.GuildMember(m.GuildID, m.Author.ID)

	for _, roleID := range member.Roles {
		for _, role := range guildRoles {
			if role.ID == roleID {
				if role.Permissions&discordgo.PermissionAdministrator != 0 {
					return true
				}
			}
		}
	}
	return false
}
