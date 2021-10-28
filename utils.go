package main

import "github.com/bwmarrin/discordgo"

func searchChannel(s *discordgo.Session, guildID string, channelName string, channelParent string) (channel *discordgo.Channel) {
	channels, _ := s.GuildChannels(guildID)

	for _, channel := range channels {
		if channel.Name == channelName {
			if channelParent == "" || channel.ParentID == channelParent {
				return channel
			}
		}
	}

	return nil
}

func searchRole(s *discordgo.Session, guildID string, roleName string) (role *discordgo.Role) {
	roles, _ := s.GuildRoles(guildID)

	for _, role := range roles {
		if role.Name == roleName {
			return role
		}
	}

	return nil
}

func searchUser(s *discordgo.Session, guildID string, userName string) (role *discordgo.User) {
	users, _ := s.GuildMembers(guildID, "", 1000)

	for _, user := range users {
		if user.User.Username == userName {
			return user.User
		}
	}

	return nil
}
