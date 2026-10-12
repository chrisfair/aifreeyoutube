package youtube

import (
	"testing"
)

// TestExtractVideos verifies that the recursive Innertube parser accurately traverses
// deeply nested JSON structures and correctly extracts videoRenderers and modern lockupViewModels.
func TestExtractVideos(t *testing.T) {
	parser := NewResponseParser()

	rawInnertubeJSON := map[string]any{
		"contents": map[string]any{
			"twoColumnSearchResultsRenderer": map[string]any{
				"primaryContents": map[string]any{
					"sectionListRenderer": map[string]any{
						"contents": []any{
							map[string]any{
								"itemSectionRenderer": map[string]any{
									"contents": []any{
										// 1. Standard videoRenderer
										map[string]any{
											"videoRenderer": map[string]any{
												"videoId": "vid_classic_01",
												"title": map[string]any{
													"runs": []any{
														map[string]any{"text": "Classic Human Woodworking Guide"},
													},
												},
												"ownerText": map[string]any{
													"runs": []any{
														map[string]any{
															"text": "Handmade Craftsmanship",
															"navigationEndpoint": map[string]any{
																"browseEndpoint": map[string]any{
																	"browseId":         "UCwoodwork123",
																	"canonicalBaseUrl": "/@woodworkcraft",
																},
															},
														},
													},
												},
												"viewCountText": map[string]any{
													"simpleText": "125,000 views",
												},
												"publishedTimeText": map[string]any{
													"simpleText": "2 months ago",
												},
												"lengthText": map[string]any{
													"simpleText": "18:42",
												},
											},
										},
										// 2. Modern lockupViewModel (Watch Next sidebar format)
										map[string]any{
											"lockupViewModel": map[string]any{
												"contentId":   "vid_lockup_02",
												"contentType": "LOCKUP_CONTENT_TYPE_VIDEO",
												"metadata": map[string]any{
													"lockupMetadataViewModel": map[string]any{
														"title": map[string]any{
															"content": "Deep In The Forge: Blacksmithing 101",
														},
														"metadata": map[string]any{
															"contentMetadataViewModel": map[string]any{
																"metadataRows": []any{
																	map[string]any{
																		"metadataParts": []any{
																			map[string]any{
																				"text": map[string]any{
																					"content": "Iron & Fire Co.",
																				},
																			},
																		},
																	},
																	map[string]any{
																		"metadataParts": []any{
																			map[string]any{
																				"text":               map[string]any{"content": "54K"},
																				"accessibilityLabel": "54,200 views",
																			},
																			map[string]any{
																				"text":               map[string]any{"content": "3w ago"},
																				"accessibilityLabel": "3 weeks ago",
																			},
																		},
																	},
																},
															},
														},
													},
												},
												"contentImage": map[string]any{
													"thumbnailViewModel": map[string]any{
														"image": map[string]any{
															"sources": []any{
																map[string]any{
																	"url": "https://i.ytimg.com/vi/vid_lockup_02/hqdefault.jpg",
																},
															},
														},
														"overlays": []any{
															map[string]any{
																"thumbnailBottomOverlayViewModel": map[string]any{
																	"badges": []any{
																		map[string]any{
																			"thumbnailBadgeViewModel": map[string]any{
																				"text": "12:15",
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	videos := parser.ExtractVideos(rawInnertubeJSON)
	if len(videos) != 2 {
		t.Fatalf("expected 2 extracted videos, got %d", len(videos))
	}

	// Verify Video 1 (videoRenderer)
	v1 := videos[0]
	if v1.ID != "vid_classic_01" {
		t.Errorf("expected video ID vid_classic_01, got %s", v1.ID)
	}
	if v1.Title != "Classic Human Woodworking Guide" {
		t.Errorf("expected title 'Classic Human Woodworking Guide', got '%s'", v1.Title)
	}
	if v1.ChannelName != "Handmade Craftsmanship" {
		t.Errorf("expected channel name 'Handmade Craftsmanship', got '%s'", v1.ChannelName)
	}
	if v1.ChannelHandle != "@woodworkcraft" {
		t.Errorf("expected channel handle '@woodworkcraft', got '%s'", v1.ChannelHandle)
	}
	if v1.Duration != "18:42" {
		t.Errorf("expected duration '18:42', got '%s'", v1.Duration)
	}
	if v1.IsShort {
		t.Errorf("expected regular video, got IsShort=true")
	}

	// Verify Video 2 (lockupViewModel)
	v2 := videos[1]
	if v2.ID != "vid_lockup_02" {
		t.Errorf("expected video ID vid_lockup_02, got %s", v2.ID)
	}
	if v2.Title != "Deep In The Forge: Blacksmithing 101" {
		t.Errorf("expected title 'Deep In The Forge: Blacksmithing 101', got '%s'", v2.Title)
	}
	if v2.ChannelName != "Iron & Fire Co." {
		t.Errorf("expected channel name 'Iron & Fire Co.', got '%s'", v2.ChannelName)
	}
	if v2.Views != "54,200 views" {
		t.Errorf("expected views '54,200 views', got '%s'", v2.Views)
	}
	if v2.Published != "3 weeks ago" {
		t.Errorf("expected published '3 weeks ago', got '%s'", v2.Published)
	}
	if v2.Duration != "12:15" {
		t.Errorf("expected duration '12:15', got '%s'", v2.Duration)
	}
}
