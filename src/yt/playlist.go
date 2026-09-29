package yt

import (
	"context"
	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func GetPlaylist(playlist_id string) ([]string, []string, error) {
	ctx := context.Background()
	var vide_ids []string
	var vide_names []string

	config, err := get_oauth_config()
	if err != nil {
		return nil, nil, err
	}

	httpClient, err := GetClient(config)
	if err != nil {
		return nil, nil, err
	}

	service, err := youtube.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, nil, err
	}

	seen := make(map[string]struct{})
	page_token := ""
	for {
		call := service.PlaylistItems.List([]string{"snippet", "contentDetails", "status"}).
			PlaylistId(playlist_id).
			MaxResults(50)
		if page_token != "" {
			call = call.PageToken(page_token)
		}

		response, err := call.Do()
		if err != nil {
			return nil, nil, check_token_error(config, err)
		}

		for _, item := range response.Items {
			if !is_downloadable(item) {
				log.Printf("Skip unavailable video: %s", video_label(item))
				continue
			}
			id := item.ContentDetails.VideoId
			// 同じ動画がプレイリストに複数あっても1回だけ扱う
			if _, dup := seen[id]; dup {
				continue
			}
			seen[id] = struct{}{}
			vide_ids = append(vide_ids, id)
			vide_names = append(vide_names, item.Snippet.Title)
		}

		if response.NextPageToken == "" {
			break
		}
		page_token = response.NextPageToken
	}

	return vide_ids, vide_names, nil
}

// is_downloadable は削除済み・非公開の動画を除く。
// これらは yt-dlp で必ず失敗し、毎回の同期で再試行されて「変更なし」にならないため
func is_downloadable(item *youtube.PlaylistItem) bool {
	if item.ContentDetails == nil || item.ContentDetails.VideoId == "" {
		return false
	}
	if item.Status != nil {
		switch item.Status.PrivacyStatus {
		case "private", "privacyStatusUnspecified":
			return false
		}
	}
	if item.Snippet != nil {
		switch item.Snippet.Title {
		case "Deleted video", "Private video":
			return false
		}
	}
	return true
}

func video_label(item *youtube.PlaylistItem) string {
	id, title := "?", ""
	if item.ContentDetails != nil && item.ContentDetails.VideoId != "" {
		id = item.ContentDetails.VideoId
	}
	if item.Snippet != nil {
		title = item.Snippet.Title
	}
	return id + " " + title
}
