package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"net/url"
)

const (
	API_URL_IP            = "https://api.ipify.org?format=json"
	API_URL_DivisionsList = "https://www.geoguessr.com/api/v4/ranked-system/divisions"
	API_Cookie_Name       = "_ncfa"
)

type IpApiResponse struct {
	IP string `json:"ip"`
}

type DivisionsListAPIResponse struct {
	Divisions []DivisionInfoAPIResponse `json:"divisions"`
}

type DivisionInfoAPIResponse struct {
	// only potentially useful fields
	DivisionNumber int    `json:"divisionNumber"`
	DivisionRank   int    `json:"divisionRank"`
	Tier           string `json:"tier"`
	Name           string `json:"name"`
	// not relying on `gameModes` as this is formatted differently
	Maps map[string]TargetMapInAPIResponse `json:"maps"`
}

type TargetMapInAPIResponse struct {
	MapID   string `json:"mapId"`
	MapName string `json:"mapName"`
}

type DivisionModeMapInfo struct {
	DivisionName string
	GameMode     string
	MapID        string
	MapName      string
}

func (dmmi DivisionModeMapInfo) String() string {
	divisionName := fmt.Sprintf("%-*s", 15, strings.ReplaceAll(dmmi.DivisionName, " ", "_"))
	modeName := fmt.Sprintf("%-*s", 8, strings.ReplaceAll(dmmi.GameMode, "Duels", ""))
	return fmt.Sprintf("DivisionModeMapInfo{DivisionName=%s  GameMode=%s MapID=%s MapName=%q}", divisionName, modeName, dmmi.MapID, dmmi.MapName)

}

func checkIP(useProxy bool) {
	var httpProxyURL string
	if !useProxy {
		httpProxyURL = ""
	} else {
		httpProxyURL = os.Getenv("HTTP_PROXY_URL")
	}

	client, err := newHTTPClient(httpProxyURL)
	if err != nil {
		panic(err)
	}

	ipReq, ipErr := http.NewRequestWithContext(context.Background(), http.MethodGet, API_URL_IP, nil)
	if ipErr != nil {
		panic(ipErr)
	}

	ipResp, ipErr := client.Do(ipReq)
	if ipErr != nil {
		panic(ipErr)
	}
	defer ipResp.Body.Close()

	ipBody, ipErr := io.ReadAll(ipResp.Body)
	if ipErr != nil {
		panic(ipErr)
	}

	if ipResp.StatusCode < 200 || ipResp.StatusCode >= 300 {
		panic(fmt.Sprintf("unexpected HTTP status: %s\nbody: %s", ipResp.Status, string(ipBody)))
	}

	var ipApiResponse IpApiResponse
	if err := json.Unmarshal(ipBody, &ipApiResponse); err != nil {
		panic(err)
	}

	if useProxy {
		fmt.Printf("Proxy IP: %s\n", ipApiResponse.IP)
	} else {
		fmt.Printf("Real IP:  %s\n", ipApiResponse.IP)
	}
}

func main() {
	cookieValue := mustEnv("GGWD_NCFA")
	httpProxyURL := os.Getenv("HTTP_PROXY_URL")

	checkIP(false)
	checkIP(true)

	client, err := newHTTPClient(httpProxyURL)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, API_URL_DivisionsList, nil)
	if err != nil {
		panic(err)
	}

	req.AddCookie(&http.Cookie{
		Name:  API_Cookie_Name,
		Value: cookieValue,
	})

	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		panic(fmt.Sprintf("unexpected HTTP status: %s\nbody: %s", resp.Status, string(body)))
	}

	rawDivisions, err := decodeObjects(body)
	if err != nil {
		panic(err)
	}

	divisionModeMapInfos := flattenObjects(rawDivisions)

	for _, dmmi := range divisionModeMapInfos {
		fmt.Println(dmmi)
	}
}

func newHTTPClient(httpProxyURL string) (*http.Client, error) {

	transport := http.DefaultTransport.(*http.Transport).Clone()
	if httpProxyURL != "" {
		parsedProxyURL, err := url.Parse(httpProxyURL)
		if err != nil {
			return nil, err
		}
		transport.Proxy = http.ProxyURL(parsedProxyURL)
	}
	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}, nil

}

func decodeObjects(body []byte) ([]DivisionInfoAPIResponse, error) {
	var wrapped DivisionsListAPIResponse
	if err := json.Unmarshal(body, &wrapped); err == nil && wrapped.Divisions != nil {
		return wrapped.Divisions, nil
	}

	return nil, errors.New("object with divisions array")
}

func flattenObjects(raw_divisions []DivisionInfoAPIResponse) []DivisionModeMapInfo {
	out := make([]DivisionModeMapInfo, 0)

	for _, divisionInfo := range raw_divisions {
		for gameMode, mapPointer := range divisionInfo.Maps {
			out = append(out, DivisionModeMapInfo{
				DivisionName: divisionInfo.Name,
				GameMode:     gameMode,
				MapID:        mapPointer.MapID,
				MapName:      mapPointer.MapName,
			})
		}
	}

	return out
}

func mustEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		panic("missing required env var: " + name)
	}
	return value
}
