package discord

import (
	"acTrackerBot/tracker"
	"acTrackerBot/tracker/acdb"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func HandlerReady(s *discordgo.Session, r *discordgo.Ready) {
	log.Println("bot is ready")
}

func HandlerCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// all bot commands start with a '!'
	if !strings.HasPrefix(m.Content, "!") {
		log.Printf("not a bot command")
		return
	}

	if strings.HasPrefix(m.Content, "!info ") {
		reg := strings.Split(m.Content, " ")
		data := acdb.GetAircraftData(reg[1])

		if data.Reg == "" {
			s.ChannelMessageSend(m.ChannelID,
				fmt.Sprintf("**%v**\nno information found, maybe registration is not valid",
					reg[1]))
		} else {
			s.ChannelMessageSend(m.ChannelID,
				fmt.Sprintf("**%v**\n```Type: %v Model: %v\nManufacturer: %v\nYear: %v\nOwner: %v```",
					data.Reg, data.Icaotype, data.Model, data.Manufacturer, data.Year, data.Ownop))
		}

	} else if strings.HasPrefix(m.Content, "!add ") {
		reg := strings.Split(m.Content, " ")
		tracker.AddRegistrationChannel <- reg[1]
	} else if strings.HasPrefix(m.Content, "!list") {
		s.ChannelMessageSend(m.ChannelID, tracker.GetRegList())
	} else if strings.HasPrefix(m.Content, "!remove ") {
		reg := strings.Split(m.Content, " ")
		tracker.RemoveRegistrationChannel <- reg[1]
	} else if strings.HasPrefix(m.Content, "!help") {
		s.ChannelMessageSend(m.ChannelID, "commands !add <reg>, !remove <reg>, !list, !info <reg>, !help")
	} else {
		s.ChannelMessageSend(m.ChannelID, "buuh, unknown command ...")
	}
}
