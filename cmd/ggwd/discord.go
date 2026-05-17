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

func NiceDivisionModeString(divisionName string, gameMode string) string {
	mode := gameMode
	if mode == "standardDuels" {
		mode = "Moving"
	} else if mode == "noMoveDuels" {
		mode = "NoMove"
	} else if mode == "nmpzDuels" {
		mode = "NMPZ"
	}
	return fmt.Sprintf("%s - %s", divisionName, mode)
}
func formatMapUrl(mapID string, mapName string) string {
	return fmt.Sprintf("[%s](https://www.geoguessr.com/maps/%s)", mapName, mapID)
}

func SendDiscordMapChangeNotification(ctx context.Context, cfg config.DiscordConfig, divisionMode DivisionModeMapDetails, delta DivisionMapDelta) error {
	notificationTitle := fmt.Sprintf("Map change in %s!", NiceDivisionModeString(divisionMode.divisionName, divisionMode.gameMode))
	notificationMessage := "No previous map information available"

	if delta.MapPointerChanged {
		notificationTitle = fmt.Sprintf("New map assigned for %s", NiceDivisionModeString(divisionMode.divisionName, divisionMode.gameMode))
		notificationMessage = formatMapUrl(divisionMode.mapInfo.ID, divisionMode.mapInfo.Name)
	} else if delta.MapDetailsChanged {
		notificationTitle = fmt.Sprintf("Map was updated for %s", NiceDivisionModeString(divisionMode.divisionName, divisionMode.gameMode))
		notificationMessage = formatMapUrl(divisionMode.mapInfo.ID, divisionMode.mapInfo.Name) + "\n\n" + delta.Details
	}

	return SendDiscordNotification(ctx, cfg, notificationTitle, notificationMessage)
}
