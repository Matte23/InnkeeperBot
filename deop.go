package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func commandDeop(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Parameters should contains the channel name and the username
	parameters := strings.Fields(m.Content)
	if len(parameters) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Please specify the channel name and the username")
		return
	}

	username := strings.TrimPrefix(m.Content, parameters[0]+" "+parameters[1]+" ")

	category := searchChannel(s, m.GuildID, "Personalizzato", "")

	if category == nil {
		s.ChannelMessageSend(m.ChannelID, "Unable to delete channel, category not found")
	}

	channel := searchChannel(s, m.GuildID, parameters[1], category.ID)
	if channel == nil {
		s.ChannelMessageSend(m.ChannelID, "Channel "+parameters[1]+" not found, unable to remove permission")
		return
	}

	user := searchUser(s, m.GuildID, username)
	if user == nil {
		s.ChannelMessageSend(m.ChannelID, "User "+username+" not found")
		return
	}

	if len(channel.PermissionOverwrites) == 0 {
		s.ChannelMessageSend(m.ChannelID, "You cannot change permissions of a public channel")
		return
	}

	// Check if the user has got permissions
	authorized := false
	for _, role := range channel.PermissionOverwrites {
		if role.ID == m.Author.ID && role.Allow&0x00000010 == 0x00000010 {
			authorized = true
		}
	}

	if authorized {
		s.ChannelPermissionSet(channel.ID, user.ID, discordgo.PermissionOverwriteTypeMember, 0x00100000, 0)
		s.ChannelMessageSend(m.ChannelID, "User "+username+" cannot manage "+parameters[1]+" anymore")
	} else {
		s.ChannelMessageSend(m.ChannelID, "You don't have the permission to revoke the authorization of an user to manage "+parameters[1])
	}
}
