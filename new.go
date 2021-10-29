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
