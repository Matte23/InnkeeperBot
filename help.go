package main

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func commandHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Description: "InnkeeperBot allow users to manage custom channels",
		Color:       0x00ff00,
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:  "!new <channel> [private]",
				Value: "Create a new channel. Add the \"private\" keyword at the end of the command to make that channel private",
			},
			&discordgo.MessageEmbedField{
				Name:  "!del <channel>",
				Value: "Delete an existing channel",
			},
			&discordgo.MessageEmbedField{
				Name:  "!add <channel> <user>",
				Value: "Give an user the permission to join a channel",
			},
			&discordgo.MessageEmbedField{
				Name:  "!rem <channel> <user>",
				Value: "Remove join permission from an user",
			},
			&discordgo.MessageEmbedField{
				Name:  "!op <channel> <user>",
				Value: "Give an user the permission to join, edit and delete a channel",
			},
			&discordgo.MessageEmbedField{
				Name:  "!deop <channel> <user>",
				Value: "Remove join/edit/delete permission from an user",
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Title:     "InnkeeperBot help",
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
