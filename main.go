package main

import (
	"acTrackerBot/config"
	"acTrackerBot/discord"
	"acTrackerBot/tracker"
	"acTrackerBot/tracker/types"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

func init() {
	config.ReadConfig()
}

func main() {

	// start the discord bot
	log.Printf("starting discord bot ...\n")

	dg, err := discordgo.New("Bot " + config.Conf.Discordbottoken)
	if err != nil {
		log.Fatal("error creating discord session,", err)
		return
	}

	dg.AddHandler(discord.HandlerCreate)
	dg.AddHandler(discord.HandlerReady)

	// open websocket for listening
	err = dg.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
		return
	}
	defer dg.Close()

	// startup aircraft tracker
	tracker.AddRegistrationChannel = make(chan string)
	tracker.RemoveRegistrationChannel = make(chan string)
	statusUpdate := tracker.StartUp()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	for {
		select {
		case acStatus := <-statusUpdate:
			processStatusUpdate(dg, acStatus)
		case <-stop:
			log.Println("Graceful shutdown")
			close(tracker.AddRegistrationChannel)
			close(tracker.RemoveRegistrationChannel)
			return
		}
	}

}

func processStatusUpdate(dg *discordgo.Session, data types.AircraftInformation) {
	log.Printf("-> rcv aircraft from tracker %+v\n", data)

	title := fmt.Sprintf("new status for %v", data.Reg)
	embed := discord.GetEmbedMessage(data)

	if err := discord.SendComplexMessageWithWebhook(dg, title, embed); err != nil {
		log.Printf("can not send complex message %v", err)
	}

}
