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
	log.Debugf("Preparing guild %s for automatic channel management", guildID)

	category := searchChannel(s, guildID, categoryName, "")
	if category == nil {
		log.Errorf("Cannot find category %s in guild %s", categoryName, guildID)
		return
	}

	deleteAllChannelsUnderCategory(s, guildID, category.ID)
	createNewChannelEndpoint(s, guildID, category.ID)

	createdChannels[guildID] = make(map[string][]string)
}

func createNewChannelEndpoint(s *discordgo.Session, guildID string, channelParent string) {

	newChannelData := discordgo.GuildChannelCreateData{
		Name:     "Crea nuovo canale",
		Type:     discordgo.ChannelTypeGuildVoice,
		ParentID: channelParent}

	var err error
	newChannelEndpoint[guildID], err = s.GuildChannelCreateComplex(guildID, newChannelData)

	if err != nil {
		log.Errorf("Cannot create channel with name %s in guild %s. This guild is in a corrupted state", newChannelData.Name, guildID)
		return
	}

	log.Debugf("Created channel with name %s in guild %s", newChannelData.Name, guildID)
}

func createNewChannel(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	category := searchChannel(s, vs.GuildID, categoryName, "")
	if category == nil {
		log.Errorf("Cannot find category %s in guild %s", categoryName, vs.GuildID)
		return
	}

	channelName := getActivity(s.State, vs.GuildID, vs.UserID)
	editChannelData := discordgo.ChannelEdit{Name: channelName}
	_, err := s.ChannelEdit(vs.ChannelID, &editChannelData)
	if err != nil {
		log.Errorf("Cannot rename channel %s with new name %s in guild %s", vs.ChannelID, channelName, vs.GuildID)
		return
	}
	createdChannels[vs.GuildID][vs.ChannelID] = []string{vs.UserID}
	log.Debugf("Renamed channel %s to %s in guild %s", vs.ChannelID, channelName, vs.GuildID)

	createNewChannelEndpoint(s, vs.GuildID, category.ID)
}

func removeUserFromChannels(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	for channelID, channel := range createdChannels[vs.GuildID] {
		index := contains(channel, vs.UserID)

		if index != -1 {
			if len(channel) == 1 {
				_, err := s.ChannelDelete(channelID)
				if err != nil {
					log.Errorf("Cannot delete channel %s in guild %s", channelID, vs.GuildID)
				} else {
					log.Debugf("Channel %s deleted in guild %s", channelID, vs.GuildID)
					delete(createdChannels[vs.GuildID], channelID)
				}
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

	channel, err := s.Channel(channelID)
	if err != nil {
		log.Errorf("Channel %s not found in guild %s", channelID, guildID)
		return
	}

	oldName := channel.Name

	if oldName == name {
		return
	}

	editChannelData := discordgo.ChannelEdit{Name: name}
	_, err = s.ChannelEdit(channelID, &editChannelData)
	if err != nil {
		log.Errorf("Cannot rename channel %s with new name %s in guild %s", channelID, name, guildID)
	}
	log.Debugf("Renamed channel %s from %s to %s in guild %s", channelID, oldName, name, guildID)
}

func getActivity(st *discordgo.State, guildID string, userID string) string {
	presence, err := st.Presence(guildID, userID)

	if err != nil || len(presence.Activities) == 0 {
		return "Stanza"
	}

	return presence.Activities[len(presence.Activities)-1].Name
}
