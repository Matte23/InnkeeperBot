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
	"github.com/bwmarrin/discordgo"
)

var newChannelEndpoint = make(map[string]*discordgo.Channel)
var createdChannels = make(map[string]map[string][]string)

var categoryName = "Gaming🎮"

func initChannels(s *discordgo.Session, guildID string) {
	category := searchChannel(s, guildID, categoryName, "")
	if category == nil {
		return
	}

	deleteAllChannelsUnderCategory(s, guildID, category.ID)
	createNewChannelEndpoint(s, guildID)

	createdChannels[guildID] = make(map[string][]string)
}

func createNewChannelEndpoint(s *discordgo.Session, guildID string) {
	category := searchChannel(s, guildID, categoryName, "")
	if category == nil {
		return
	}

	newChannelData := discordgo.GuildChannelCreateData{
		Name:     "Crea nuovo canale",
		Type:     discordgo.ChannelTypeGuildVoice,
		ParentID: category.ID}

	newChannelEndpoint[guildID], _ = s.GuildChannelCreateComplex(guildID, newChannelData)
}

func createNewChannel(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	channelName := getActivity(s.State, vs.GuildID, vs.UserID)
	s.ChannelEdit(vs.ChannelID, channelName)
	createdChannels[vs.GuildID][vs.ChannelID] = []string{vs.UserID}

	createNewChannelEndpoint(s, vs.GuildID)
}

func removeUserFromChannels(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	for channelID, channel := range createdChannels[vs.GuildID] {
		index := contains(channel, vs.UserID)

		if index != -1 {
			if len(channel) == 1 {
				s.ChannelDelete(channelID)
			} else {
				createdChannels[vs.GuildID][channelID] = remove(channel, index)
			}
		}
	}
}

func updateChannelName(s *discordgo.Session, guildID string, channelID string) {
	programs := make(map[string]int)

	for _, userID := range createdChannels[guildID][channelID] {
		programs[getActivity(s.State, guildID, userID)] += 1
	}

	max := -1
	name := "Stanza"

	for program, count := range programs {
		if count > max {
			max = count
			name = program
		}
	}

	s.ChannelEdit(channelID, name)
}

func getActivity(st *discordgo.State, guildID string, userID string) string {
	presence, err := st.Presence(guildID, userID)

	if err != nil || len(presence.Activities) == 0 {
		return "Stanza"
	}

	return presence.Activities[len(presence.Activities)-1].Name
}
