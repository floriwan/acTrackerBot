package discord

import (
	"acTrackerBot/tracker"
	"acTrackerBot/tracker/aeroDataBox"
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
		if reg[1] == "" {
			fmt.Printf("error: no registration set to request aircarft data")
			return
		}

		data, err := aeroDataBox.GetAircraftInfo(reg[1])

		if err != nil {
			_, err := s.ChannelMessageSend(m.ChannelID,
				fmt.Sprintf("**%v**\nno information found, maybe registration is not valid",
					reg[1]))
			if err != nil {
				log.Printf("error %v\n", err)
			}
			return
		}

		embed := GetEmbedAircraftMessage(*data)
		if err := SendComplexMessageWithWebhook(s,
			fmt.Sprintf("aircraft information for %v", data.Reg), embed); err != nil {
			fmt.Printf("error, can not send embed aircraft information %+v %v\n", &data, err)
		}

		/*
			data := acdb.GetAircraftData(reg[1])

			if data.Reg == "" {
				_, err := s.ChannelMessageSend(m.ChannelID,
					fmt.Sprintf("**%v**\nno information found, maybe registration is not valid",
						reg[1]))
				if err != nil {
					log.Printf("error %v\n", err)
				}
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
		*/
	} else if strings.HasPrefix(m.Content, "!add ") {
		reg := strings.Split(m.Content, " ")
		if reg[1] == "" {
			fmt.Printf("error: no registration set to add to tracker")
			return
		}
		tracker.AddRegistrationChannel <- reg[1]
	} else if strings.HasPrefix(m.Content, "!list") {
		_, err := s.ChannelMessageSend(m.ChannelID, tracker.GetRegList())
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
	} else if strings.HasPrefix(m.Content, "!remove ") {
		reg := strings.Split(m.Content, " ")
		if reg[1] == "" {
			fmt.Printf("error: no registration set to remove from tracker")
			return
		}
		tracker.RemoveRegistrationChannel <- reg[1]
	} else if strings.HasPrefix(m.Content, "!help") {
		_, err := s.ChannelMessageSend(m.ChannelID, "commands !add <reg>, !remove <reg>, !list, !info <reg>, !help")
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
	} else {
		_, err := s.ChannelMessageSend(m.ChannelID, "buuh, unknown command ...")
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
	}
}
