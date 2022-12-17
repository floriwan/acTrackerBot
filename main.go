package main

import (
	"acTrackerBot/config"
	"acTrackerBot/discord"
	"acTrackerBot/tracker"
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
			log.Printf("-> rcv aircraft from tracker %v\n", acStatus)

			title := fmt.Sprintf("new status for %v", acStatus.Reg)
			embed := []*discordgo.MessageEmbed{{
				Fields: []*discordgo.MessageEmbedField{
					{Name: "Departure", Value: acStatus.Origin, Inline: true},
					{Name: "Arrival", Value: acStatus.Destination, Inline: true},
					{Name: "Status", Value: acStatus.Status.String()},
					{Name: "Speed (kt)", Value: fmt.Sprintf("%v", acStatus.Speed), Inline: true},
					{Name: "Altitude (ft)", Value: fmt.Sprintf("%v", acStatus.AltGeom), Inline: true},
					{Name: "Squawk", Value: acStatus.Squawk},
					{Name: "Lat", Value: fmt.Sprintf("%v", acStatus.Lat), Inline: true},
					{Name: "Long", Value: fmt.Sprintf("%v", acStatus.Lon), Inline: true},
				},
			}}

			discord.SendComplexMessageWithWebhook(dg, title, embed)

			/*discord.SendDiscordMessageWithWebhook(dg,
			fmt.Sprintf("Registration: **%v** Callsign: **%v**\n```%v->%v status: %v\nground speed: %v kt\nalt geom: %v ft\nlat: %v long: %v```",
				acStatus.Reg, acStatus.Callsign,
				acStatus.Origin, acStatus.Destination,
				acStatus.Status.String(), acStatus.Speed, acStatus.AltGeom, acStatus.Lat, acStatus.Lon))
			*/
		case <-stop:
			log.Println("Graceful shutdown")
			close(tracker.AddRegistrationChannel)
			close(tracker.RemoveRegistrationChannel)
			return
		}
	}

}
