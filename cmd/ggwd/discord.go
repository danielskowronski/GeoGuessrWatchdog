package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/types"
)

func SendDiscordNotification(ctx context.Context, urls []string, title string, message string, color string) error {
	if len(urls) == 0 {
		return nil
	}

	for i, urlStr := range urls {
		u, err := url.Parse(urlStr)
		if err != nil {
			return fmt.Errorf("parse shoutrrr URL: %w", err)
		}
		if u.Scheme == "discord" { // this hack is needed because of https://github.com/containrrr/shoutrrr/issues/396
			q := u.Query()
			q.Set("splitLines", "false")
			u.RawQuery = q.Encode()
			urls[i] = u.String()
		}
	}
	sender, err := shoutrrr.CreateSender(urls...)
	if err != nil {
		return fmt.Errorf("create shoutrrr sender: %w", err)
	}

	if color == "" {
		color = "#FFFFFF"
	}

	params := &types.Params{
		"title":      title,
		"color":      color,
		"splitlines": "false", // this doesn't work directly
	}

	errorsList := sender.Send(message, params)
	for _, err := range errorsList {
		if err != nil {
			return fmt.Errorf("send shoutrrr notification: %w", err)
		}
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

func SendDiscordMapChangeNotification(ctx context.Context, urls []string, divisionMode DivisionModeMapDetails, delta DivisionMapDelta) error {
	notificationTitle := fmt.Sprintf("Map change in %s!", NiceDivisionModeString(divisionMode.divisionName, divisionMode.gameMode))
	notificationMessage := "No previous map information available"
	color := "#800080"

	if delta.MapPointerChanged {
		notificationTitle = fmt.Sprintf("New map assigned for %s", NiceDivisionModeString(divisionMode.divisionName, divisionMode.gameMode))
		notificationMessage = formatMapUrl(divisionMode.mapInfo.ID, divisionMode.mapInfo.Name)
		color = "#FF0000"
	} else if delta.MapDetailsChanged {
		notificationTitle = fmt.Sprintf("Map was updated for %s", NiceDivisionModeString(divisionMode.divisionName, divisionMode.gameMode))
		notificationMessage = formatMapUrl(divisionMode.mapInfo.ID, divisionMode.mapInfo.Name) + "\n\n" + delta.Details
	}

	return SendDiscordNotification(ctx, urls, notificationTitle, notificationMessage, color)
}
