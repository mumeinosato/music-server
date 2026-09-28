package yt

import (
	"context"

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

	page_token := ""
	for {
		call := service.PlaylistItems.List([]string{"snippet", "contentDetails"}).
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
			vide_ids = append(vide_ids, item.ContentDetails.VideoId)
			vide_names = append(vide_names, item.Snippet.Title)
		}

		if response.NextPageToken == "" {
			break
		}
		page_token = response.NextPageToken
	}

	return vide_ids, vide_names, nil
}
