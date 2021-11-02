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
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(handleMessage)
	dg.AddHandler(handleVoiceActivity)
	dg.AddHandler(handlePresenceUpdate)
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func handleVoiceActivity(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	if val, ok := newChannelEndpoint[vs.GuildID]; ok {
		if vs.ChannelID == "" {
			removeUserFromChannels(s, vs)
			return
		}

		if vs.ChannelID == val.ID {
			removeUserFromChannels(s, vs)
			createNewChannel(s, vs)
		} else if _, ok := createdChannels[vs.GuildID][vs.ChannelID]; ok {
			createdChannels[vs.GuildID][vs.ChannelID] = append(createdChannels[vs.GuildID][vs.ChannelID], vs.UserID)
		} else {
			removeUserFromChannels(s, vs)
		}
	} else {
		initChannels(s, vs.GuildID)
	}
}

func handlePresenceUpdate(s *discordgo.Session, pu *discordgo.PresenceUpdate) {
	if _, ok := newChannelEndpoint[pu.GuildID]; ok {
		for channelID, channel := range createdChannels[pu.GuildID] {
			for _, user := range channel {
				if user == pu.User.ID {
					updateChannelName(s, pu.GuildID, channelID)
				}
			}
		}
	} else {
		initChannels(s, pu.GuildID)
	}
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message starts with "!new", create a new custom channel. Parameters: name, privacy [public or private](default public, optional)
	if strings.HasPrefix(m.Content, "!new") {
		commandNew(s, m)
	}

	// If the message starts with "!del", delete a custom channel. Parameters: name
	if strings.HasPrefix(m.Content, "!del") {
		commandDelete(s, m)
	}

	// If the message starts with "!add", enable an user to join the channel. Parameters: channel name, user name
	if strings.HasPrefix(m.Content, "!add") {
		commandAdd(s, m)
	}

	// If the message starts with "!rem", remove from an user the permission to join the channel. Parameters: channel name, user name
	if strings.HasPrefix(m.Content, "!rem") {
		commandRem(s, m)
	}

	// If the message starts with "!op", enable an user to manage the channel. Parameters: channel name, user name
	if strings.HasPrefix(m.Content, "!op") {
		commandOp(s, m)
	}

	// If the message starts with "!deop", remove from an user the permission to manage the channel. Parameters: channel name, user name
	if strings.HasPrefix(m.Content, "!deop") {
		commandDeop(s, m)
	}

	// If the message starts with "!rpc", remove from an user the permission to join the channel. Parameters: channel name, user name
	if strings.HasPrefix(m.Content, "!help") {
		commandHelp(s, m)
	}
}
