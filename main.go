package main

import (
	"acTrackerBot/config"
	"acTrackerBot/tracker"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	config.ReadConfig()
}

func main() {
	fmt.Printf("starting discord bot ...\n")

	dg, err := discordgo.New("Bot " + config.Conf.Discordbottoken)
	if err != nil {
		fmt.Println("error creating discord session,", err)
		return
	}

	dg.AddHandler(handlerReady)
	dg.AddHandler(handlerCreate)

	// open websocket for listening
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	defer dg.Close()

	// startup aircraft tracker
	ch := tracker.StartUp()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	for {
		select {
		case msg := <-ch:
			fmt.Printf("rcv msg from tracker %v\n", msg)
		case <-stop:
			log.Println("Graceful shutdown")
			return
		}
	}

}

func handlerReady(s *discordgo.Session, r *discordgo.Ready) {
	fmt.Println("bot is ready ...")
}

func handlerCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// all bot commands start with a '!'
	if !strings.HasPrefix(m.Content, "!") {
		fmt.Printf("not a bot command")
		return
	}

	if strings.HasPrefix(m.Content, "!add ") {
		reg := strings.Split(m.Content, " ")
		fmt.Printf("add new registration %v\n", reg[1])
		tracker.Add(reg[1])
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("new registration list size: %v", tracker.Size()))
	} else if strings.HasPrefix(m.Content, "!list") {
		fmt.Printf("list all registrations\n")
		s.ChannelMessageSend(m.ChannelID, tracker.List())
	} else if strings.HasPrefix(m.Content, "!remove ") {
		reg := strings.Split(m.Content, " ")
		fmt.Printf("remove registration %v\n", reg[1])
		tracker.Remove(reg[1])
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("new registration list size: %v", tracker.Size()))
	} else if strings.HasPrefix(m.Content, "!help") {
		s.ChannelMessageSend(m.ChannelID, "commands !add <reg>, !remove <reg>, !list, !help")
	} else {
		s.ChannelMessageSend(m.ChannelID, "buuh, unknown command ...")
	}
}
