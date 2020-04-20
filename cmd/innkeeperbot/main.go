// InnkeeperBot, a bot to create and manage custom channels
// Copyright (C) 2020 Matteo Schiff
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
	"time"

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
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message starts with "!new", create a new custom channel. Parameters: name, privacy [public or private](default public, optional)
	if strings.HasPrefix(m.Content, "!new") {
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
		s.ChannelPermissionSet(newChannel.ID, s.State.User.ID, "member", 0x00100000|0x00000010, 0)
		s.ChannelPermissionSet(newChannel.ID, m.Author.ID, "member", 0x00100000|0x00000010, 0)

		if len(parameters) > 2 && parameters[2] == "private" {
			roleEveryone := searchRole(s, m.GuildID, "@everyone")
			// Deny anyone else joining this channel
			s.ChannelPermissionSet(newChannel.ID, roleEveryone.ID, "role", 0, 0x00100000)
		}

		s.ChannelMessageSend(m.ChannelID, "New channel created: "+parameters[1])
	}

	// If the message starts with "!del", delete a custom channel. Parameters: name
	if strings.HasPrefix(m.Content, "!del") {
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

	// If the message starts with "!add", enable an user to join the channel. Parameters: channel name, user name
	if strings.HasPrefix(m.Content, "!add") {
		// Parameters should contains the channel name and the username
		parameters := strings.Fields(m.Content)
		if len(parameters) < 2 {
			s.ChannelMessageSend(m.ChannelID, "Please specify the channel name")
			return
		}

		category := searchChannel(s, m.GuildID, "Personalizzato", "")
		if category == nil {
			s.ChannelMessageSend(m.ChannelID, "Unable to delete channel, category not found")
		}

		channel := searchChannel(s, m.GuildID, parameters[1], category.ID)
		if channel == nil {
			s.ChannelMessageSend(m.ChannelID, "Channel "+parameters[1]+" not found, unable to add permission")
			return
		}

		user := searchUser(s, m.GuildID, parameters[2])
		if user == nil {
			s.ChannelMessageSend(m.ChannelID, "User "+parameters[2]+" not found")
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
			s.ChannelPermissionSet(channel.ID, user.ID, "member", 0x00100000, 0)
			s.ChannelMessageSend(m.ChannelID, "User "+parameters[2]+" can now join "+parameters[1])
		} else {
			s.ChannelMessageSend(m.ChannelID, "You don't have the permission to authorize an user to join "+parameters[1])
		}
	}

	// If the message starts with "!rem", remove from an user the permission to join the channel. Parameters: channel name, user name
	if strings.HasPrefix(m.Content, "!rem") {
		// Parameters should contains the channel name and the username
		parameters := strings.Fields(m.Content)
		if len(parameters) < 2 {
			s.ChannelMessageSend(m.ChannelID, "Please specify the channel name")
			return
		}

		category := searchChannel(s, m.GuildID, "Personalizzato", "")

		if category == nil {
			s.ChannelMessageSend(m.ChannelID, "Unable to delete channel, category not found")
		}

		channel := searchChannel(s, m.GuildID, parameters[1], category.ID)
		if channel == nil {
			s.ChannelMessageSend(m.ChannelID, "Channel "+parameters[1]+" not found, unable to remove permission")
			return
		}

		user := searchUser(s, m.GuildID, parameters[2])
		if user == nil {
			s.ChannelMessageSend(m.ChannelID, "User "+parameters[2]+" not found")
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
			s.ChannelMessageSend(m.ChannelID, "User "+parameters[2]+" cannot join "+parameters[1]+" anymore")
		} else {
			s.ChannelMessageSend(m.ChannelID, "You don't have the permission to revoke the authorization of an user to join "+parameters[1])
		}

	}

	// If the message starts with "!op", enable an user to manage the channel. Parameters: channel name, user name
	if strings.HasPrefix(m.Content, "!op") {
		// Parameters should contains the channel name and the username
		parameters := strings.Fields(m.Content)
		if len(parameters) < 2 {
			s.ChannelMessageSend(m.ChannelID, "Please specify the channel name")
			return
		}

		category := searchChannel(s, m.GuildID, "Personalizzato", "")
		if category == nil {
			s.ChannelMessageSend(m.ChannelID, "Unable to delete channel, category not found")
		}

		channel := searchChannel(s, m.GuildID, parameters[1], category.ID)
		if channel == nil {
			s.ChannelMessageSend(m.ChannelID, "Channel "+parameters[1]+" not found, unable to add permission")
			return
		}

		user := searchUser(s, m.GuildID, parameters[2])
		if user == nil {
			s.ChannelMessageSend(m.ChannelID, "User "+parameters[2]+" not found")
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
			s.ChannelPermissionSet(channel.ID, user.ID, "member", 0x00100000|0x00000010, 0)
			s.ChannelMessageSend(m.ChannelID, "User "+parameters[2]+" can now manage "+parameters[1])
		} else {
			s.ChannelMessageSend(m.ChannelID, "You don't have the permission to authorize an user to manage "+parameters[1])
		}
	}

	// If the message starts with "!rem", remove from an user the permission to manage the channel. Parameters: channel name, user name
	if strings.HasPrefix(m.Content, "!deop") {
		// Parameters should contains the channel name and the username
		parameters := strings.Fields(m.Content)
		if len(parameters) < 2 {
			s.ChannelMessageSend(m.ChannelID, "Please specify the channel name")
			return
		}

		category := searchChannel(s, m.GuildID, "Personalizzato", "")

		if category == nil {
			s.ChannelMessageSend(m.ChannelID, "Unable to delete channel, category not found")
		}

		channel := searchChannel(s, m.GuildID, parameters[1], category.ID)
		if channel == nil {
			s.ChannelMessageSend(m.ChannelID, "Channel "+parameters[1]+" not found, unable to remove permission")
			return
		}

		user := searchUser(s, m.GuildID, parameters[2])
		if user == nil {
			s.ChannelMessageSend(m.ChannelID, "User "+parameters[2]+" not found")
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
			s.ChannelPermissionSet(channel.ID, user.ID, "member", 0x00100000, 0)
			s.ChannelMessageSend(m.ChannelID, "User "+parameters[2]+" cannot manage "+parameters[1]+" anymore")
		} else {
			s.ChannelMessageSend(m.ChannelID, "You don't have the permission to revoke the authorization of an user to manage "+parameters[1])
		}

	}

	// If the message starts with "!rpc", remove from an user the permission to join the channel. Parameters: channel name, user name
	if strings.HasPrefix(m.Content, "!help") {

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
}

func searchChannel(s *discordgo.Session, guildID string, channelName string, channelParent string) (channel *discordgo.Channel) {
	channels, _ := s.GuildChannels(guildID)

	for _, channel := range channels {
		if channel.Name == channelName {
			if channelParent == "" || channel.ParentID == channelParent {
				return channel
			}
		}
	}

	return nil
}

func searchRole(s *discordgo.Session, guildID string, roleName string) (role *discordgo.Role) {
	roles, _ := s.GuildRoles(guildID)

	for _, role := range roles {
		if role.Name == roleName {
			return role
		}
	}

	return nil
}

func searchUser(s *discordgo.Session, guildID string, userName string) (role *discordgo.User) {
	users, _ := s.GuildMembers(guildID, "", 1000)

	for _, user := range users {
		if user.User.Username == userName {
			return user.User
		}
	}

	return nil
}
