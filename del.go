package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func commandDelete(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Parameters should contains the name of the channel to delete
	parameters := strings.Fields(m.Content)
	if len(parameters) < 1 {
		s.ChannelMessageSend(m.ChannelID, "Please specify the channel name")
		return
	}

	category := searchChannel(s, m.GuildID, "Personalizzato", "")
	if category == nil {
		s.ChannelMessageSend(m.ChannelID, "Unable to delete channel, category not found")
		return
	}

	channel := searchChannel(s, m.GuildID, parameters[1], category.ID)
	if channel == nil {
		s.ChannelMessageSend(m.ChannelID, "Unable to delete channel, not found")
	} else if len(channel.PermissionOverwrites) == 0 {
		s.ChannelDelete(channel.ID)
		s.ChannelMessageSend(m.ChannelID, "Channel "+parameters[1]+" deleted")
	} else {
		// Check if the user has got manage channels permissions
		authorized := false
		for _, role := range channel.PermissionOverwrites {
			if role.ID == m.Author.ID && role.Allow&0x00000010 == 0x00000010 {
				authorized = true
			}
		}

		if authorized {
			s.ChannelDelete(channel.ID)
			s.ChannelMessageSend(m.ChannelID, "Channel "+parameters[1]+" deleted")
		} else {
			s.ChannelMessageSend(m.ChannelID, "You don't have the permission to delete this channel")
		}
	}
}
