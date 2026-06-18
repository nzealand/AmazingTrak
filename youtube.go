package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var youtubeHTTPClient = &http.Client{Timeout: 4 * time.Second}

// fetchYouTubeTitle retrieves a video's title via YouTube's public oEmbed
// endpoint (no API key required). Best-effort: any failure returns ("", false)
// and callers should fall back to whatever title the user supplied.
func fetchYouTubeTitle(videoURL string) (string, bool) {
	oembedURL := "https://www.youtube.com/oembed?format=json&url=" + url.QueryEscape(videoURL)
	resp, err := youtubeHTTPClient.Get(oembedURL)
	if err != nil {
		return "", false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", false
	}
	var data struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", false
	}
	title := strings.TrimSpace(data.Title)
	if title == "" {
		return "", false
	}
	return title, true
}
