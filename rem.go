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

func commandRem(s *discordgo.Session, m *discordgo.MessageCreate) {
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
		s.ChannelPermissionDelete(channel.ID, user.ID)
		s.ChannelMessageSend(m.ChannelID, "User "+username+" cannot join "+parameters[1]+" anymore")
	} else {
		s.ChannelMessageSend(m.ChannelID, "You don't have the permission to revoke the authorization of an user to join "+parameters[1])
	}
}
