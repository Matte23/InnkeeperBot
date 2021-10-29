package main

import (
	"github.com/bwmarrin/discordgo"
)

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

func getAllChannelsUnderCategory(s *discordgo.Session, guildID string, channelParent string) (channels []*discordgo.Channel) {
	allChannels, _ := s.GuildChannels(guildID)

	for _, channel := range allChannels {
		if channelParent == "" || channel.ParentID == channelParent {
			channels = append(channels, channel)
		}
	}

	return channels
}

func deleteAllChannelsUnderCategory(s *discordgo.Session, guildID string, channelParent string) {
	for _, channel := range getAllChannelsUnderCategory(s, guildID, channelParent) {
		s.ChannelDelete(channel.ID)
	}
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

func contains(s []string, e string) int {
	for i, a := range s {
		if a == e {
			return i
		}
	}
	return -1
}

func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
