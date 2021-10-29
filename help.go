// InnkeeperBot, a bot to create and manage custom channels
// Copyright (C) 2020-2021 Matteo Schiff
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
