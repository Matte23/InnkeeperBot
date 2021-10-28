package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func commandNew(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Parameters should contains the channel name and the privacy preference
	parameters := strings.Fields(m.Content)
	if len(parameters) < 1 {
		s.ChannelMessageSend(m.ChannelID, "Please specify the channel name")
		return
	}

	category := searchChannel(s, m.GuildID, "Personalizzato", "")
	if category == nil {
		s.ChannelMessageSend(m.ChannelID, "Unable to create channel, category not found")
		return
	}

	newChannelData := discordgo.GuildChannelCreateData{
		Name:     parameters[1],
		Type:     discordgo.ChannelTypeGuildVoice,
		ParentID: category.ID}

	newChannel, _ := s.GuildChannelCreateComplex(m.GuildID, newChannelData)

	// Allow the channel creator and the bot to join the channel
	s.ChannelPermissionSet(newChannel.ID, s.State.User.ID, discordgo.PermissionOverwriteTypeMember, 0x00100000|0x00000010, 0)
	s.ChannelPermissionSet(newChannel.ID, m.Author.ID, discordgo.PermissionOverwriteTypeMember, 0x00100000|0x00000010, 0)

	if len(parameters) > 2 && parameters[2] == "private" {
		roleEveryone := searchRole(s, m.GuildID, "@everyone")
		// Deny anyone else joining this channel
		s.ChannelPermissionSet(newChannel.ID, roleEveryone.ID, discordgo.PermissionOverwriteTypeRole, 0, 0x00100000)
	}

	s.ChannelMessageSend(m.ChannelID, "New channel created: "+parameters[1])
}
