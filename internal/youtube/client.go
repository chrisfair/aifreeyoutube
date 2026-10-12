package youtube

import (
	"context"
	"fmt"
	"strings"

	"aifreeyoutube/internal/models"
)

// Service provides high-level video discovery and inspection methods backed by Innertube.
type Service struct {
	innertube *InnertubeClient
	parser    *ResponseParser
}

// NewService instantiates a YouTube discovery service.
func NewService() *Service {
	return &Service{
		innertube: NewInnertubeClient(),
		parser:    NewResponseParser(),
	}
}

// Search queries YouTube's internal search endpoint for the provided query string.
// It parses matching video items, formats timestamps, thumbnails, and author information.
func (s *Service) Search(ctx context.Context, query string) ([]models.Video, error) {
	if strings.TrimSpace(query) == "" {
		return s.BrowseHome(ctx)
	}

	payload := SearchPayload{
		Context: s.innertube.buildContext(),
		Query:   query,
	}

	rawResp, err := s.innertube.executePostRequest(ctx, "search", payload)
	if err != nil {
		return nil, fmt.Errorf("search failed for query '%s': %w", query, err)
	}

	videos := s.parser.ExtractVideos(rawResp)
	return videos, nil
}

// BrowseHome fetches the default clean discovery home feed.
func (s *Service) BrowseHome(ctx context.Context) ([]models.Video, error) {
	return s.Search(ctx, "popular documentary science craftsmanship")
}

// BrowseTrending queries YouTube's trending discovery feed.
func (s *Service) BrowseTrending(ctx context.Context) ([]models.Video, error) {
	return s.Search(ctx, "trending documentary science technology craftsmanship")
}

// GetVideoDetails retrieves the metadata for a specific video ID alongside its related recommendations.
func (s *Service) GetVideoDetails(ctx context.Context, videoID string) (models.Video, []models.Video, error) {
	payload := NextPayload{
		Context: s.innertube.buildContext(),
		VideoID: videoID,
	}

	rawResp, err := s.innertube.executePostRequest(ctx, "next", payload)
	if err != nil {
		return models.Video{}, nil, fmt.Errorf("failed to fetch video details for %s: %w", videoID, err)
	}

	primaryVideo, relatedVideos := s.parser.ExtractWatchPage(rawResp, videoID)
	return primaryVideo, relatedVideos, nil
}
