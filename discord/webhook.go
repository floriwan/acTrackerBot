package discord

import (
	"acTrackerBot/config"

	"github.com/bwmarrin/discordgo"
)

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
