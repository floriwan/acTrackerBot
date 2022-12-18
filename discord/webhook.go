package discord

import (
	"acTrackerBot/config"
	"acTrackerBot/tracker/aeroDataBox"
	"acTrackerBot/tracker/types"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func GetEmbedAircraftMessage(data aeroDataBox.Aircraft) (embed []*discordgo.MessageEmbed) {

	fields := []*discordgo.MessageEmbedField{}
	if data.AirlineName != "" {
		fields = append(fields, greateEmbededField("Airline", data.AirlineName, false))
	}

	if data.IataCodeShort != "" {
		fields = append(fields, greateEmbededField("IATA Code", data.IataCodeShort, true))
	}

	if data.IcaoCode != "" {
		fields = append(fields, greateEmbededField("ICAO Code", data.IcaoCode, true))
	}

	if data.ProductionLine != "" {
		fields = append(fields, greateEmbededField("Production Line", data.ProductionLine, true))
	}

	if data.TypeName != "" {
		fields = append(fields, greateEmbededField("Type", data.TypeName, true))
	}

	if data.Model != "" {
		fields = append(fields, greateEmbededField("Model", data.Model, true))
	}

	if data.ModelCode != "" {
		fields = append(fields, greateEmbededField("Model Code", data.ModelCode, true))
	}

	if data.EngineType != "" {
		fields = append(fields, greateEmbededField("Engine Type", data.EngineType, true))
	}

	if data.FirstFlightDate != "" {
		fields = append(fields, greateEmbededField("First Flight", data.FirstFlightDate, true))
	}

	if data.RegistrationDate != "" {
		fields = append(fields, greateEmbededField("Registration Date", data.RegistrationDate, true))
	}

	fields = append(fields, greateEmbededField("Number Of Seats", fmt.Sprintf("%v", data.NumSeats), true))
	fields = append(fields, greateEmbededField("Number Of Engines", fmt.Sprintf("%v", data.NumEngines), true))

	return []*discordgo.MessageEmbed{{
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: data.Image.Url,
		},
		Fields: fields,
	}}
}

func greateEmbededField(name string, data string, inline bool) (field *discordgo.MessageEmbedField) {
	field = &discordgo.MessageEmbedField{
		Name:   name,
		Value:  data,
		Inline: inline,
	}
	return field
}

func GetEmbedMessage(data types.AircraftInformation) (embed []*discordgo.MessageEmbed) {

	squawk := "unknown"

	if data.Squawk != "" {
		squawk = data.Squawk
	}

	return []*discordgo.MessageEmbed{{
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Flight", Value: fmt.Sprintf("%v > %v", data.Origin, data.Destination), Inline: true},
			{Name: "Flight Status", Value: data.FlightStatus, Inline: true},
			{Name: "Flight Callsign", Value: data.Callsign, Inline: true},
			{Name: "Status", Value: data.Status.String()},
			{Name: "Speed (kt)", Value: fmt.Sprintf("%v", data.Speed), Inline: true},
			{Name: "Altitude (ft)", Value: fmt.Sprintf("%v", data.AltGeom), Inline: true},
			{Name: "Squawk", Value: squawk, Inline: true},
			{Name: "Lat", Value: fmt.Sprintf("%v", data.Lat), Inline: true},
			{Name: "Long", Value: fmt.Sprintf("%v", data.Lon), Inline: true},
		},
	}}
}

func SendComplexMessageWithWebhook(dg *discordgo.Session,
	content string,
	embeds []*discordgo.MessageEmbed) error {

	webhook, err := dg.WebhookWithToken(config.Conf.Discrodwebhookid, config.Conf.Discrodwebhooktoken)
	if err != nil {
		return err
	}

	_, err = dg.ChannelMessageSendComplex(webhook.ChannelID, &discordgo.MessageSend{
		Content: content,
		Embeds:  embeds,
	})

	if err != nil {
		return err
	}

	return nil
}

func SendDiscordMessageWithWebhook(dg *discordgo.Session, msg string) error {
	webhook, err := dg.WebhookWithToken(config.Conf.Discrodwebhookid, config.Conf.Discrodwebhooktoken)
	if err != nil {
		return err
	}

	_, err = dg.ChannelMessageSend(webhook.ChannelID, msg)
	if err != nil {
		return err
	}
	return nil
}
