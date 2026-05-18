package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/types"
)

const (
	SHOUTRRR_URI_KEY_PREPEND = "ggwd_prepend"
	SHOUTRRR_URI_KEY_APPEND  = "ggwd_append"
)

type ShoutrrrExtraOptions struct {
	Append  string
	Prepend string
}

func ParseShoutrrrUri(rawURL string) (string, ShoutrrrExtraOptions, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", ShoutrrrExtraOptions{}, fmt.Errorf("parse shoutrrr url: %s", err)
	}
	extraOptions := ShoutrrrExtraOptions{}
	query := u.Query()
	extraOptions.Append = query.Get(SHOUTRRR_URI_KEY_APPEND)
	extraOptions.Prepend = query.Get(SHOUTRRR_URI_KEY_PREPEND)
	query.Del(SHOUTRRR_URI_KEY_APPEND)
	query.Del(SHOUTRRR_URI_KEY_PREPEND)
	u.RawQuery = query.Encode()
	return u.String(), extraOptions, nil
}
func PatchShoutrrrUriToSetSplitLinesFalse(rawURL string) (string, error) {
	// hack needed until https://github.com/containrrr/shoutrrr/pull/498 is merged
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL, fmt.Errorf("parse shoutrrr URL: %w", err)
	}
	if u.Scheme == "discord" {
		q := u.Query()
		q.Set("splitLines", "false")
		u.RawQuery = q.Encode()
		return u.String(), nil
	}
	return rawURL, nil
}

func SendDiscordNotification(ctx context.Context, urls []string, title string, message string, color string) error {
	if len(urls) == 0 {
		return nil
	}
	errorsList := make([]error, 0)

	for _, urlRaw := range urls {
		urlRaw, err := PatchShoutrrrUriToSetSplitLinesFalse(urlRaw)
		if err != nil {
			return fmt.Errorf("patch shoutrrr URL: %w", err)
		}
		url, extraOptions, err := ParseShoutrrrUri(urlRaw)
		if err != nil {
			return fmt.Errorf("parse shoutrrr URL: %w", err)
		}
		if extraOptions.Prepend != "" {
			message = extraOptions.Prepend + message
		}
		if extraOptions.Append != "" {
			message = message + extraOptions.Append
		}

		sender, err := shoutrrr.CreateSender([]string{url}...)
		if err != nil {
			return fmt.Errorf("create shoutrrr sender: %w", err)
		}
		if color == "" {
			color = "#FFFFFF"
		}
		params := &types.Params{
			"title":      title,
			"color":      color,
			"splitLines": "false",
		}
		errorsList = append(errorsList, sender.Send(message, params)...)
	}

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
