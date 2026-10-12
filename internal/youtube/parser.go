package youtube

import (
	"fmt"
	"strings"

	"aifreeyoutube/internal/models"
)

// ResponseParser provides methods to traverse complex, nested Innertube response objects
// and map them into standardized models.Video structures.
type ResponseParser struct{}

// NewResponseParser creates a parser instance.
func NewResponseParser() *ResponseParser {
	return &ResponseParser{}
}

// ExtractVideos searches an arbitrary Innertube JSON dictionary recursively for video objects.
// Innertube nests video items deeply within sectionListRenderers, itemSectionRenderers,
// videoRenderers, and compactVideoRenderers. This recursive walker guarantees robust extraction
// regardless of UI variations YouTube delivers across regions or experiment buckets.
func (p *ResponseParser) ExtractVideos(root map[string]any) []models.Video {
	var videos []models.Video
	seenIDs := make(map[string]bool)

	var walk func(node any)
	walk = func(node any) {
		switch val := node.(type) {
		case map[string]any:
			// 1. Standard search or browse feed item
			if vr, ok := val["videoRenderer"].(map[string]any); ok {
				if v, valid := p.parseVideoRenderer(vr); valid && !seenIDs[v.ID] {
					seenIDs[v.ID] = true
					videos = append(videos, v)
				}
			}

			// 2. Compact video item (found on Watch Page related recommendations)
			if cvr, ok := val["compactVideoRenderer"].(map[string]any); ok {
				if v, valid := p.parseCompactVideoRenderer(cvr); valid && !seenIDs[v.ID] {
					seenIDs[v.ID] = true
					videos = append(videos, v)
				}
			}

			// 3. Grid video item (found on topic pages or channel videos)
			if gvr, ok := val["gridVideoRenderer"].(map[string]any); ok {
				if v, valid := p.parseGridVideoRenderer(gvr); valid && !seenIDs[v.ID] {
					seenIDs[v.ID] = true
					videos = append(videos, v)
				}
			}

			// 4. Shorts reel item (can be extracted and marked as Short)
			if rir, ok := val["reelItemRenderer"].(map[string]any); ok {
				if v, valid := p.parseReelItemRenderer(rir); valid && !seenIDs[v.ID] {
					seenIDs[v.ID] = true
					videos = append(videos, v)
				}
			}

			// 5. Modern Lockup View Model (modern YouTube Innertube recommendation format)
			// YouTube's updated desktop WEB client replaces compactVideoRenderer in the Watch
			// Page secondary recommendations with lockupViewModel structures.
			if lvm, ok := val["lockupViewModel"].(map[string]any); ok {
				if v, valid := p.parseLockupViewModel(lvm); valid && !seenIDs[v.ID] {
					seenIDs[v.ID] = true
					videos = append(videos, v)
				}
			}

			// Continue recursive descent through map children
			for _, child := range val {
				walk(child)
			}

		case []any:
			// Continue recursive descent through array children
			for _, item := range val {
				walk(item)
			}
		}
	}

	walk(root)
	return videos
}

// parseVideoRenderer converts an Innertube "videoRenderer" structure into a models.Video.
func (p *ResponseParser) parseVideoRenderer(vr map[string]any) (models.Video, bool) {
	videoID, _ := vr["videoId"].(string)
	if videoID == "" {
		return models.Video{}, false
	}

	title := p.extractText(vr["title"])
	channelName, channelID, channelHandle := p.extractOwnerInfo(vr["ownerText"], vr["longBylineText"])
	views := p.extractText(vr["viewCountText"])
	published := p.extractText(vr["publishedTimeText"])
	duration := p.extractText(vr["lengthText"])
	thumbnail := p.extractThumbnail(vr["thumbnail"], videoID)
	description := p.extractDescription(vr["detailedMetadataSnippets"], vr["descriptionSnippet"])

	isShort := strings.Contains(duration, ":") == false && (strings.Contains(title, "#shorts") || strings.Contains(title, "#short"))
	if duration != "" && p.isDurationUnderMinute(duration) {
		isShort = true
	}

	return models.Video{
		ID:            videoID,
		Title:         title,
		ChannelName:   channelName,
		ChannelID:     channelID,
		ChannelHandle: channelHandle,
		Views:         views,
		Published:     published,
		Duration:      duration,
		Description:   description,
		Thumbnail:     thumbnail,
		EmbedURL:      fmt.Sprintf("https://www.youtube-nocookie.com/embed/%s?autoplay=1&rel=0&modestbranding=1", videoID),
		IsShort:       isShort,
	}, true
}

// parseCompactVideoRenderer handles related videos found on the watch page sidebar.
func (p *ResponseParser) parseCompactVideoRenderer(cvr map[string]any) (models.Video, bool) {
	videoID, _ := cvr["videoId"].(string)
	if videoID == "" {
		return models.Video{}, false
	}

	title := p.extractText(cvr["title"])
	channelName, channelID, channelHandle := p.extractOwnerInfo(cvr["shortBylineText"], cvr["longBylineText"])
	views := p.extractText(cvr["viewCountText"])
	published := p.extractText(cvr["publishedTimeText"])
	duration := p.extractText(cvr["lengthText"])
	thumbnail := p.extractThumbnail(cvr["thumbnail"], videoID)

	return models.Video{
		ID:            videoID,
		Title:         title,
		ChannelName:   channelName,
		ChannelID:     channelID,
		ChannelHandle: channelHandle,
		Views:         views,
		Published:     published,
		Duration:      duration,
		Thumbnail:     thumbnail,
		EmbedURL:      fmt.Sprintf("https://www.youtube-nocookie.com/embed/%s?autoplay=1&rel=0&modestbranding=1", videoID),
	}, true
}

// parseGridVideoRenderer handles channel uploads and browse grid layouts.
func (p *ResponseParser) parseGridVideoRenderer(gvr map[string]any) (models.Video, bool) {
	videoID, _ := gvr["videoId"].(string)
	if videoID == "" {
		return models.Video{}, false
	}

	title := p.extractText(gvr["title"])
	channelName, channelID, channelHandle := p.extractOwnerInfo(gvr["shortBylineText"], nil)
	views := p.extractText(gvr["viewCountText"])
	published := p.extractText(gvr["publishedTimeText"])
	thumbnail := p.extractThumbnail(gvr["thumbnail"], videoID)

	return models.Video{
		ID:            videoID,
		Title:         title,
		ChannelName:   channelName,
		ChannelID:     channelID,
		ChannelHandle: channelHandle,
		Views:         views,
		Published:     published,
		Thumbnail:     thumbnail,
		EmbedURL:      fmt.Sprintf("https://www.youtube-nocookie.com/embed/%s?autoplay=1&rel=0&modestbranding=1", videoID),
	}, true
}

// parseReelItemRenderer parses YouTube Shorts entries.
func (p *ResponseParser) parseReelItemRenderer(rir map[string]any) (models.Video, bool) {
	videoID, _ := rir["videoId"].(string)
	if videoID == "" {
		return models.Video{}, false
	}

	title := p.extractText(rir["headline"])
	views := p.extractText(rir["viewCountText"])
	thumbnail := p.extractThumbnail(rir["thumbnail"], videoID)

	return models.Video{
		ID:        videoID,
		Title:     title,
		Views:     views,
		Duration:  "Short",
		Thumbnail: thumbnail,
		EmbedURL:  fmt.Sprintf("https://www.youtube-nocookie.com/embed/%s?autoplay=1&rel=0&modestbranding=1", videoID),
		IsShort:   true,
	}, true
}

// parseLockupViewModel translates YouTube Innertube's modern "lockupViewModel" component into a models.Video.
//
// Technical Context:
// In YouTube's 2024+ client schema rollout, the legacy `compactVideoRenderer` used on the watch page sidebar
// (watch next recommendations) has been superseded by `lockupViewModel`. Unlike the legacy protobuf-derived
// renderers that used nested `runs` and `simpleText` containers, ViewModels encapsulate their visual hierarchy
// under structured metadata payloads:
//   - Video ID: Provided in `contentId`
//   - Content Type: `contentType` (e.g., "LOCKUP_CONTENT_TYPE_VIDEO")
//   - Title: Located in `metadata.lockupMetadataViewModel.title.content`
//   - Metadata Rows: Contained in `metadata.lockupMetadataViewModel.metadata.contentMetadataViewModel.metadataRows`
//       * Row 0: Creator/Channel display name
//       * Row 1: View count and relative upload timestamp
//   - Thumbnail & Duration: Located under `contentImage.thumbnailViewModel`
//       * Image sources: High/low resolution URLs in `image.sources`
//       * Duration Badge: Stored in `overlays[].thumbnailBottomOverlayViewModel.badges[].thumbnailBadgeViewModel.text`
//
// This method handles full graceful degradation if any nested sub-objects are absent or formatted differently.
func (p *ResponseParser) parseLockupViewModel(lvm map[string]any) (models.Video, bool) {
	videoID, _ := lvm["contentId"].(string)
	if videoID == "" {
		return models.Video{}, false
	}

	// Filter out non-video content if contentType is explicitly specified and not a video
	if cType, ok := lvm["contentType"].(string); ok && cType != "" && cType != "LOCKUP_CONTENT_TYPE_VIDEO" {
		return models.Video{}, false
	}

	var title string
	var channelName string
	var views string
	var published string
	var duration string
	var thumbnail string

	// Extract Title & Text Metadata
	if metaContainer, ok := lvm["metadata"].(map[string]any); ok {
		if lmvm, ok := metaContainer["lockupMetadataViewModel"].(map[string]any); ok {
			// Video title
			if tObj, ok := lmvm["title"].(map[string]any); ok {
				if c, ok := tObj["content"].(string); ok {
					title = strings.TrimSpace(c)
				}
			}

			// Channel, views, and timestamp rows
			if cmvmContainer, ok := lmvm["metadata"].(map[string]any); ok {
				if cmvm, ok := cmvmContainer["contentMetadataViewModel"].(map[string]any); ok {
					if rows, ok := cmvm["metadataRows"].([]any); ok {
						// Row 0 typically contains the Creator / Channel name
						if len(rows) > 0 {
							if r0, ok := rows[0].(map[string]any); ok {
								if parts, ok := r0["metadataParts"].([]any); ok && len(parts) > 0 {
									if p0, ok := parts[0].(map[string]any); ok {
										if t, ok := p0["text"].(map[string]any); ok {
											channelName, _ = t["content"].(string)
										}
									}
								}
							}
						}

						// Row 1 typically contains view count and publication date
						if len(rows) > 1 {
							if r1, ok := rows[1].(map[string]any); ok {
								if parts, ok := r1["metadataParts"].([]any); ok {
									if len(parts) > 0 {
										if p0, ok := parts[0].(map[string]any); ok {
											if aLabel, ok := p0["accessibilityLabel"].(string); ok && aLabel != "" {
												views = aLabel
											} else if t, ok := p0["text"].(map[string]any); ok {
												views, _ = t["content"].(string)
											}
										}
									}
									if len(parts) > 1 {
										if p1, ok := parts[1].(map[string]any); ok {
											if aLabel, ok := p1["accessibilityLabel"].(string); ok && aLabel != "" {
												published = aLabel
											} else if t, ok := p1["text"].(map[string]any); ok {
												published, _ = t["content"].(string)
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// Extract Thumbnail and Duration from contentImage
	if contentImage, ok := lvm["contentImage"].(map[string]any); ok {
		if tvm, ok := contentImage["thumbnailViewModel"].(map[string]any); ok {
			// Thumbnail URL (take highest resolution source available)
			if img, ok := tvm["image"].(map[string]any); ok {
				if sources, ok := img["sources"].([]any); ok && len(sources) > 0 {
					lastSrc := sources[len(sources)-1]
					if sMap, ok := lastSrc.(map[string]any); ok {
						thumbnail, _ = sMap["url"].(string)
					}
				}
			}

			// Overlay Duration Badge (e.g. "14:59" or "3:21")
			if overlays, ok := tvm["overlays"].([]any); ok {
				for _, o := range overlays {
					if oMap, ok := o.(map[string]any); ok {
						if bOverlay, ok := oMap["thumbnailBottomOverlayViewModel"].(map[string]any); ok {
							if badges, ok := bOverlay["badges"].([]any); ok && len(badges) > 0 {
								if b0, ok := badges[0].(map[string]any); ok {
									if tbvm, ok := b0["thumbnailBadgeViewModel"].(map[string]any); ok {
										duration, _ = tbvm["text"].(string)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// Graceful defaults
	if thumbnail == "" {
		thumbnail = fmt.Sprintf("https://i.ytimg.com/vi/%s/hqdefault.jpg", videoID)
	}
	if channelName == "" {
		channelName = "YouTube Creator"
	}

	isShort := strings.Contains(duration, ":") == false && (strings.Contains(title, "#shorts") || strings.Contains(title, "#short"))
	if duration != "" && p.isDurationUnderMinute(duration) {
		isShort = true
	}

	return models.Video{
		ID:          videoID,
		Title:       title,
		ChannelName: channelName,
		Views:       views,
		Published:   published,
		Duration:    duration,
		Thumbnail:   thumbnail,
		EmbedURL:    fmt.Sprintf("https://www.youtube-nocookie.com/embed/%s?autoplay=1&rel=0&modestbranding=1", videoID),
		IsShort:     isShort,
	}, true
}

// extractText inspects Innertube text containers (which can be simpleText or an array of runs).
func (p *ResponseParser) extractText(container any) string {
	if container == nil {
		return ""
	}
	m, ok := container.(map[string]any)
	if !ok {
		return ""
	}

	// Case 1: Simple text field
	if simple, ok := m["simpleText"].(string); ok && simple != "" {
		return strings.TrimSpace(simple)
	}

	// Case 2: Array of "runs"
	if runs, ok := m["runs"].([]any); ok {
		var sb strings.Builder
		for _, r := range runs {
			if runMap, ok := r.(map[string]any); ok {
				if t, ok := runMap["text"].(string); ok {
					sb.WriteString(t)
				}
			}
		}
		return strings.TrimSpace(sb.String())
	}

	return ""
}

// extractOwnerInfo parses channel display name, browse ID, and handle from owner/byline structures.
func (p *ResponseParser) extractOwnerInfo(primaryByline, secondaryByline any) (name, id, handle string) {
	target := primaryByline
	if target == nil {
		target = secondaryByline
	}
	if target == nil {
		return "Unknown Creator", "", ""
	}

	m, ok := target.(map[string]any)
	if !ok {
		return "Unknown Creator", "", ""
	}

	name = p.extractText(m)

	if runs, ok := m["runs"].([]any); ok && len(runs) > 0 {
		if firstRun, ok := runs[0].(map[string]any); ok {
			if nav, ok := firstRun["navigationEndpoint"].(map[string]any); ok {
				if browseEndpoint, ok := nav["browseEndpoint"].(map[string]any); ok {
					if browseID, ok := browseEndpoint["browseId"].(string); ok {
						id = browseID
					}
					if canonicalURL, ok := browseEndpoint["canonicalBaseUrl"].(string); ok {
						handle = strings.TrimPrefix(canonicalURL, "/")
					}
				}
			}
		}
	}

	if name == "" {
		name = "YouTube Creator"
	}
	return name, id, handle
}

// extractThumbnail retrieves the highest quality thumbnail available, or falls back to standard i.ytimg.com format.
func (p *ResponseParser) extractThumbnail(container any, videoID string) string {
	fallback := fmt.Sprintf("https://i.ytimg.com/vi/%s/hqdefault.jpg", videoID)
	if container == nil {
		return fallback
	}

	m, ok := container.(map[string]any)
	if !ok {
		return fallback
	}

	thumbnails, ok := m["thumbnails"].([]any)
	if !ok || len(thumbnails) == 0 {
		return fallback
	}

	// Innertube orders thumbnails ascending by resolution. The last item is the highest resolution.
	lastItem := thumbnails[len(thumbnails)-1]
	if lastMap, ok := lastItem.(map[string]any); ok {
		if u, ok := lastMap["url"].(string); ok && u != "" {
			return u
		}
	}

	return fallback
}

// extractDescription parses search snippet descriptions and detail fragments.
func (p *ResponseParser) extractDescription(snippets, descriptionSnippet any) string {
	var sb strings.Builder

	if snippets != nil {
		if arr, ok := snippets.([]any); ok {
			for _, s := range arr {
				if smap, ok := s.(map[string]any); ok {
					if snippetText, ok := smap["snippetText"]; ok {
						sb.WriteString(p.extractText(snippetText))
						sb.WriteString(" ")
					}
				}
			}
		}
	}

	if sb.Len() == 0 && descriptionSnippet != nil {
		sb.WriteString(p.extractText(descriptionSnippet))
	}

	return strings.TrimSpace(sb.String())
}

// isDurationUnderMinute checks if a MM:SS length is shorter than 60 seconds (indicative of a Short).
func (p *ResponseParser) isDurationUnderMinute(duration string) bool {
	parts := strings.Split(duration, ":")
	if len(parts) == 2 && parts[0] == "0" {
		return true
	}
	return false
}

// ExtractWatchPage parses primary video metadata and secondary recommendations from the /next endpoint.
func (p *ResponseParser) ExtractWatchPage(root map[string]any, videoID string) (models.Video, []models.Video) {
	primary := models.Video{
		ID:        videoID,
		Title:     fmt.Sprintf("Video %s", videoID),
		EmbedURL:  fmt.Sprintf("https://www.youtube-nocookie.com/embed/%s?autoplay=1&rel=0&modestbranding=1", videoID),
		Thumbnail: fmt.Sprintf("https://i.ytimg.com/vi/%s/hqdefault.jpg", videoID),
	}

	var related []models.Video

	contents, ok := root["contents"].(map[string]any)
	if !ok {
		return primary, related
	}

	twoCol, ok := contents["twoColumnWatchNextResults"].(map[string]any)
	if !ok {
		return primary, related
	}

	// 1. Primary video details
	if results, ok := twoCol["results"].(map[string]any); ok {
		if innerResults, ok := results["results"].(map[string]any); ok {
			if items, ok := innerResults["contents"].([]any); ok {
				for _, itm := range items {
					if itmMap, ok := itm.(map[string]any); ok {
						if pInfo, ok := itmMap["videoPrimaryInfoRenderer"].(map[string]any); ok {
							if t := p.extractText(pInfo["title"]); t != "" {
								primary.Title = t
							}
							if vCount, ok := pInfo["viewCount"].(map[string]any); ok {
								if vrend, ok := vCount["videoViewCountRenderer"].(map[string]any); ok {
									primary.Views = p.extractText(vrend["viewCount"])
								}
							}
							if date := p.extractText(pInfo["dateText"]); date != "" {
								primary.Published = date
							}
						}

						if sInfo, ok := itmMap["videoSecondaryInfoRenderer"].(map[string]any); ok {
							if owner, ok := sInfo["owner"].(map[string]any); ok {
								if vOwner, ok := owner["videoOwnerRenderer"].(map[string]any); ok {
									name, cid, handle := p.extractOwnerInfo(vOwner["title"], nil)
									if name != "" {
										primary.ChannelName = name
									}
									primary.ChannelID = cid
									primary.ChannelHandle = handle
								}
							}
							if attrDesc, ok := sInfo["attributedDescription"].(map[string]any); ok {
								if c, ok := attrDesc["content"].(string); ok && c != "" {
									primary.Description = c
								}
							}
							if primary.Description == "" {
								primary.Description = p.extractText(sInfo["description"])
							}
						}
					}
				}
			}
		}
	}

	// 2. Related recommendations in secondaryResults
	if sec, ok := twoCol["secondaryResults"].(map[string]any); ok {
		related = p.ExtractVideos(sec)
	}

	return primary, related
}
