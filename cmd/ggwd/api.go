package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/danielskowronski/geoguessrwatchdog/internal/apischema"
	"github.com/danielskowronski/geoguessrwatchdog/internal/config"
	"go.temporal.io/sdk/activity"
)

type ApiClient struct {
	client      *http.Client
	apiBase     string
	cookieValue string

	cacheConfig       config.GeoguessrAPICacheConfig
	extraCacheTargets map[string]string
}

func NewAPIClient(httpProxyURL string, apiBase string, cookie string, cache config.GeoguessrAPICacheConfig, extraCacheTargets map[string]string) (*ApiClient, error) {
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
		apiBase:           apiBase,
		cookieValue:       cookie,
		cacheConfig:       cache,
		extraCacheTargets: extraCacheTargets,
	}, nil
}

func (a *ApiClient) GetPublicIP(ctx context.Context, ipInfoCheckURL string) (string, error) {
	body, code, err := a.FetchJSON(ctx, ipInfoCheckURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch public IP: %w", err)
	}
	if code != http.StatusOK {
		return "", fmt.Errorf("unexpected status code when fetching public IP: %d", code)
	}
	var ipInfo struct {
		IP string `json:"ip"`
	}
	if err := json.Unmarshal(body, &ipInfo); err != nil {
		return "", fmt.Errorf("failed to parse IP info response: %w", err)
	}
	return ipInfo.IP, nil
}

func (a *ApiClient) CacheResponse(targetURL string, responseBody []byte) error {
	if !a.cacheConfig.Enabled {
		return nil
	}

	var targetPath string
	if strings.HasPrefix(targetURL, a.apiBase) {
		// GeoGuessr API URLs are always cached, path is 1:1 with URL
		url := strings.ReplaceAll(targetURL, a.apiBase, "")
		targetPath = fmt.Sprintf("%s/%s", a.cacheConfig.CachePath, url)
	} else if extraPath, ok := a.extraCacheTargets[targetURL]; ok {
		// other URLs are cached only if explicitly configured, and then with a fixed path
		targetPath = fmt.Sprintf("%s/%s", a.cacheConfig.CachePath, extraPath)
	} else {
		// not a URL we want to cache
		return nil
	}

	targetDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}
	if err := os.WriteFile(targetPath, responseBody, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}
	return nil
}

func (a *ApiClient) FetchJSON(ctx context.Context, targetURL string) (json.RawMessage, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, 0, err
	}

	if strings.HasPrefix(targetURL, a.apiBase) {
		// only add cookie for GeoGuessr API requests, not for other APIs
		req.AddCookie(&http.Cookie{
			Name:  GG_COOKIE_NAME,
			Value: a.cookieValue,
		})
	}
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

	a.CacheResponse(targetURL, body)

	return raw, resp.StatusCode, nil
}

func (a *ApiClient) FetchDivisions(ctx context.Context) ([]apischema.DivisionModeMapInfo, int, error) {
	logger := activity.GetLogger(ctx)
	url := a.apiBase + GG_API_PATH_DIVISIONS_LIST
	logger.Debug("fetching divisions list from API", "url", url)
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

func (a *ApiClient) FetchUserSingleplayerStats(ctx context.Context, userID string) (*apischema.SingleplayerStatsInfo, int, error) {
	body, code, err := a.FetchJSON(ctx, a.apiBase+fmt.Sprintf(GG_API_PATH_USER_SINGLEPLAYER, userID))
	if err != nil {
		return nil, code, err
	}

	singleplayerStatsInfo, err := apischema.DecodeApiResponseSingleplayerStats(body)
	if err != nil {
		return nil, code, fmt.Errorf("failed to decode API response: %w", err)
	}

	return singleplayerStatsInfo, code, nil
}
