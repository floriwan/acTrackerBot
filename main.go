package main

import (
	"acTrackerBot/config"
	"acTrackerBot/tracker"
	"acTrackerBot/tracker/acdb"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
)

/* webhook url
https://discord.com/api/webhooks/1053637034155188345/5gnU29NBB0dhz4bsYAZNCvl9kv_J8VVu6-FhTtZWFdilGng0Q9BYY6p6XMpQCHuEd2D2
*/

func init() {
	config.ReadConfig()
}

func main() {
	log.Printf("starting discord bot ...\n")

	dg, err := discordgo.New("Bot " + config.Conf.Discordbottoken)
	if err != nil {
		log.Fatal("error creating discord session,", err)
		return
	}

	dg.AddHandler(handlerReady)
	dg.AddHandler(handlerCreate)

	// open websocket for listening
	err = dg.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
		return
	}
	defer dg.Close()

	// startup aircraft tracker
	statusUpdate := tracker.StartUp()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	for {
		select {
		case acStatus := <-statusUpdate:
			log.Printf("-> rcv aircraft from tracker %v\n", acStatus)

			webhook, err := dg.WebhookWithToken(config.Conf.Discrodwebhookid, config.Conf.Discrodwebhooktoken)
			if err != nil {
				log.Fatal("unable to open webhook", err)
			}

			dg.ChannelMessageSend(webhook.ChannelID,
				fmt.Sprintf("**%v** Type: %v\n```status: %v\nground speed: %v\nlat: %v lon: %v```",
					acStatus.IcaoType, acStatus.Reg, acStatus.Status.String(), acStatus.Speed, acStatus.Lat, acStatus.Lon))
		case <-stop:
			log.Println("Graceful shutdown")
			return
		}
	}

}

func handlerReady(s *discordgo.Session, r *discordgo.Ready) {
	log.Println("bot is ready")
}

func handlerCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
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
		if err := tracker.AddNewReg(reg[1]); err != nil {
			s.ChannelMessageSend(m.ChannelID, ("registration is not valid"))
		}
		//s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("new registration list: %v", tracker.GetRegList()))
	} else if strings.HasPrefix(m.Content, "!list") {
		s.ChannelMessageSend(m.ChannelID, tracker.GetRegList())
	} else if strings.HasPrefix(m.Content, "!remove ") {
		reg := strings.Split(m.Content, " ")
		if err := tracker.RemoveReg(reg[1]); err != nil {
			s.ChannelMessageSend(m.ChannelID, ("registration is not valid"))
		}
		//s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("new registration list: %v", tracker.GetRegList()))
	} else if strings.HasPrefix(m.Content, "!help") {
		s.ChannelMessageSend(m.ChannelID, "commands !add <reg>, !remove <reg>, !list, !info <reg>, !help")
	} else {
		s.ChannelMessageSend(m.ChannelID, "buuh, unknown command ...")
	}
}
