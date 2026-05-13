package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/danielskowronski/geoguessrwatchdog/internal/apischema"
)

type ApiClient struct {
	client      *http.Client
	apiBase     string
	cookieValue string
}

func NewAPIClient(httpProxyURL string, apiBase string, cookie string) (*ApiClient, error) {
	tr := http.DefaultTransport.(*http.Transport).Clone()
	if httpProxyURL != "" {
		u, err := url.Parse(httpProxyURL)
		if err != nil {
			return nil, err
		}
		tr.Proxy = http.ProxyURL(u)
	}
	return &ApiClient{
		client: &http.Client{
			Transport: tr,
			Timeout:   30 * time.Second, // TODO: parametrize this
		},
		apiBase:     apiBase,
		cookieValue: cookie,
	}, nil
}

func (a *ApiClient) FetchJSON(ctx context.Context, targetURL string) (json.RawMessage, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, 0, err
	}

	req.AddCookie(&http.Cookie{
		Name:  GG_COOKIE_NAME,
		Value: a.cookieValue,
	})
	// TODO: either clearly indentifying ourselves with User-Agent or spoofing real web browser

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, resp.StatusCode, fmt.Errorf("unexpected HTTP status %s: %s", resp.Status, string(body))
	}
	var raw json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, resp.StatusCode, fmt.Errorf("invalid JSON: %w", err)
	}
	return raw, resp.StatusCode, nil
}

func (a *ApiClient) FetchDivisions(ctx context.Context) ([]apischema.DivisionModeMapInfo, int, error) {
	url := a.apiBase + GG_API_PATH_DIVISIONS_LIST
	fmt.Printf("Fetching divisions list from API: %s\n", url) // FIXME: change to logger
	body, code, err := a.FetchJSON(ctx, url)
	if err != nil {
		return nil, code, err
	}

	divisions, err := apischema.DecodeApiResponseDivisionList(body)
	if err != nil {
		return nil, code, fmt.Errorf("failed to decode API response: %w", err)
	}

	return divisions, code, nil
}

func (a *ApiClient) FetchMapInfo(ctx context.Context, mapID string) (*apischema.MapInfo, int, error) {
	body, code, err := a.FetchJSON(ctx, a.apiBase+fmt.Sprintf(GG_API_PATH_MAP_INFO, mapID))
	if err != nil {
		return nil, code, err
	}

	mapInfo, err := apischema.DecodeApiResponseMap(body)
	if err != nil {
		return nil, code, fmt.Errorf("failed to decode API response: %w", err)
	}

	return mapInfo, code, nil
}

func (a *ApiClient) FetchUserProgress(ctx context.Context, userID string) (*apischema.ProgressInfo, int, error) {
	body, code, err := a.FetchJSON(ctx, a.apiBase+fmt.Sprintf(GG_API_PATH_USER_PROGRESS, userID))
	if err != nil {
		return nil, code, err
	}

	progressInfo, err := apischema.DecodeApiResponseProgress(body)
	if err != nil {
		return nil, code, fmt.Errorf("failed to decode API response: %w", err)
	}

	return progressInfo, code, nil
}

func (a *ApiClient) FetchUserStats(ctx context.Context, userID string) (*apischema.StatsInfo, int, error) {
	body, code, err := a.FetchJSON(ctx, a.apiBase+fmt.Sprintf(GG_API_PATH_USER_STATS, userID))
	if err != nil {
		return nil, code, err
	}

	statsInfo, err := apischema.DecodeApiResponseStats(body)
	if err != nil {
		return nil, code, fmt.Errorf("failed to decode API response: %w", err)
	}

	return statsInfo, code, nil
}
