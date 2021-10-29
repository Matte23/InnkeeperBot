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

func getActivity(st *discordgo.State, guildID string, userID string) string {
	presence, err := st.Presence(guildID, userID)

	if err != nil || len(presence.Activities) == 0 {
		return "Stanza"
	}

	return presence.Activities[len(presence.Activities)-1].Name
}
