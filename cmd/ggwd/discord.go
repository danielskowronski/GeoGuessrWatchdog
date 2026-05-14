package main

import (
	"context"
	"fmt"

	"github.com/danielskowronski/geoguessrwatchdog/internal/config"
	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/discord"
)

func SendDiscordNotification(ctx context.Context, cfg config.DiscordConfig, title string, message string) error {
	if cfg.BotToken == "" {
		return fmt.Errorf("discord bot token is empty")
	}
	if len(cfg.Receivers) == 0 {
		return fmt.Errorf("discord receivers list is empty")
	}

	session := discord.DefaultSession()
	discordSvc := discord.New()
	discordSvc.SetClient(session)

	if err := discordSvc.AuthenticateWithBotToken(cfg.BotToken); err != nil {
		return fmt.Errorf("authenticate discord bot: %w", err)
	}

	for _, receiver := range cfg.Receivers {
		if receiver == "" {
			return fmt.Errorf("discord receiver/channel ID is empty")
		}
		discordSvc.AddReceivers(receiver)
	}

	notifier := notify.New()
	notifier.UseServices(discordSvc)

	if err := notifier.Send(ctx, title, message); err != nil {
		return fmt.Errorf("send discord notification: %w", err)
	}

	return nil
}
