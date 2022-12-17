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
		//log.Printf("%+v\n", data)
		if data.Reg == "" {
			s.ChannelMessageSend(m.ChannelID,
				fmt.Sprintf("**%v**\nno information found, maybe registration is not valid",
					reg[1]))
		} else {

			_, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Content: fmt.Sprintf("information for %v", data.Reg),
				Embeds: []*discordgo.MessageEmbed{
					{
						//Title: data.Reg,
						//Description: fmt.Sprintf("ICAO type: %v model: %v", data.Icaotype, data.Model),
						Fields: []*discordgo.MessageEmbedField{
							{Name: "ICAO", Value: data.Icaotype, Inline: true},
							{Name: "Model", Value: data.Model, Inline: true},
							{Name: "Manufacturer", Value: data.Manufacturer},
							{Name: "Year", Value: data.Year},
							{Name: "Owner", Value: data.Ownop},
						},
					},
				},
			})
			if err != nil {
				log.Printf("error %v\n", err)
			}

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
