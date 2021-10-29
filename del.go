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
