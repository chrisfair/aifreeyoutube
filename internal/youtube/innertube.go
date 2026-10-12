// Package youtube provides a direct integration layer with YouTube's internal
// Innertube API. This allows our backend to query search results, browse trending feeds,
// and extract video metadata without requiring Google Cloud API keys, OAuth credentials,
// or being subject to restrictive daily API quotas (such as the standard 10,000 units/day).
package youtube

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// BaseInnertubeURL is the root endpoint for internal YouTube web client requests.
	BaseInnertubeURL = "https://www.youtube.com/youtubei/v1"

	// ClientVersion corresponds to a stable modern YouTube Web client build identifier.
	ClientVersion = "2.20240901.01.00"

	// ClientName specifies the client platform identifier recognized by Innertube.
	ClientName = "WEB"

	// DefaultUserAgent mimics a standard Google Chrome browser on Linux desktop to prevent bot blocking.
	DefaultUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36"
)

// InnertubeClient manages low-level HTTP communication with YouTube's internal APIs.
type InnertubeClient struct {
	httpClient *http.Client
}

// NewInnertubeClient initializes an InnertubeClient with configured HTTP connection pooling
// and sensible network timeouts.
func NewInnertubeClient() *InnertubeClient {
	return &InnertubeClient{
		httpClient: &http.Client{
			Timeout: 12 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 20,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

// ClientContext models the client identity block required in all Innertube JSON requests.
// YouTube uses this block to verify client capabilities, locale, and interface language.
type ClientContext struct {
	Client struct {
		ClientName    string `json:"clientName"`
		ClientVersion string `json:"clientVersion"`
		Hl            string `json:"hl"` // Host Language: English ("en")
		Gl            string `json:"gl"` // Geolocation: United States ("US")
	} `json:"client"`
}

// SearchPayload models the request body posted to /youtubei/v1/search.
type SearchPayload struct {
	Context ClientContext `json:"context"`
	Query   string        `json:"query"`
}

// BrowsePayload models the request body posted to /youtubei/v1/browse for feeds.
type BrowsePayload struct {
	Context  ClientContext `json:"context"`
	BrowseID string        `json:"browseId"`
}

// NextPayload models the request body posted to /youtubei/v1/next for video details and recommendations.
type NextPayload struct {
	Context ClientContext `json:"context"`
	VideoID string        `json:"videoId"`
}

// buildContext constructs the standard Innertube client context header.
func (c *InnertubeClient) buildContext() ClientContext {
	var ctx ClientContext
	ctx.Client.ClientName = ClientName
	ctx.Client.ClientVersion = ClientVersion
	ctx.Client.Hl = "en"
	ctx.Client.Gl = "US"
	return ctx
}

// executePostRequest executes an authenticated-like POST request to an Innertube endpoint
// with necessary headers and returns the unmarshaled response as a generic map.
func (c *InnertubeClient) executePostRequest(ctx context.Context, endpoint string, payload any) (map[string]any, error) {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal innertube payload: %w", err)
	}

	url := fmt.Sprintf("%s/%s?prettyPrint=false", BaseInnertubeURL, endpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", DefaultUserAgent)
	req.Header.Set("X-YouTube-Client-Name", "1") // 1 = Web Client
	req.Header.Set("X-YouTube-Client-Version", ClientVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("innertube network request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodySnippet, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("innertube returned non-200 status code: %d, body: %s", resp.StatusCode, string(bodySnippet))
	}

	var result map[string]any
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode innertube JSON response: %w", err)
	}

	return result, nil
}
