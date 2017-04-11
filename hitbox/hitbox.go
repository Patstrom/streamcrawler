package hitbox

import (
	"repos/patstrom/streamcrawler/api"
	"strconv"
)

// Hitbox implements the API interface fot the api.hitbox.tv api
type Hitbox struct{}

// Streams hold the relevant information from the top-level JSON
type Streams struct {
	Error        bool     `json:"error"` // If json key "error" exists something must've gone wrong
	ErrorMessage string   `json:"error_msg"`
	Streams      []Stream `json:"livestream"`
}

// stream holds the relevant information in the "streams" key of the original JSON
type Stream struct {
	Id              string      `json:"media_id"`
	PreviewImageUrl string      `json:"media_thumbnail_large"`
	Status          string      `json:"media_status"`
	DisplayName     string      `json:"media_display_name"`
	Viewers         string      `json:"media_views"`
	Channel         ChannelInfo `json:"channel"`
}

// channelInfo holds the relevant information in the "channel" key of the original JSON
type ChannelInfo struct {
	Id           string `json:"user_id"`
	Name         string `json:"user_name"`
	Url          string `json:"channel_link"`
	LogoImageUrl string `json:"user_logo"`
}

func (h Hitbox) Streams(game, limit string) api.Streams {
	// Define parameters to API
	gameQ := api.Pair{"game", game}
	limitQ := api.Pair{"limit", limit}
	seoQ := api.Pair{"seo", "true"}
	liveOnlyQ := api.Pair{"liveonly", "true"}

	// Hitbox doesn't require authentication
	streams := Streams{}
	success := api.ApiCall(&streams, "https://api.hitbox.tv/media/list/list/", &api.Pair{"", ""},
		gameQ, limitQ, seoQ, liveOnlyQ)

	if success {
		return convertStreamInfo(&streams)
	} else {
		return api.Streams{
			Error:        true,
			ErrorMessage: "Failed to retrieve information from API",
		}
	}

}

func convertStreamInfo(hitboxStreams *Streams) (apiStreams api.Streams) {
	apiStreams.Error = hitboxStreams.Error // Hitbox return a bool for error
	apiStreams.ErrorMessage = hitboxStreams.ErrorMessage
	apiStreams.Total = float64(len(hitboxStreams.Streams))

	for _, v := range hitboxStreams.Streams {

		id, _ := strconv.ParseFloat(v.Id, 64)
		channelId, _ := strconv.ParseFloat(v.Channel.Id, 64)
		viewers, _ := strconv.ParseFloat(v.Viewers, 64)

		apiStreams.Streams = append(apiStreams.Streams, api.Stream{
			Id:              id,
			PreviewImageUrl: v.PreviewImageUrl,
			LogoImageUrl:    v.Channel.LogoImageUrl,
			Viewers:         viewers,
			ChannelId:       channelId,
			UserName:        v.Channel.Name,
			DisplayName:     v.DisplayName,
			StatusText:      v.Status})
	}
	return
}
