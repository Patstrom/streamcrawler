package twitch

import (
	"repos/patstrom/streamcrawler/api"
)

// Twitch implements the API interface for the api.twitch.tv api
type Twitch struct {
	auth api.Pair
}

// NewTwitch constructs a new Twitch struct using the key "Client-Id" and the given string as the value.
func NewTwitch(clientID string) *Twitch {
	return &Twitch{auth: api.Pair{Key: "Client-Id", Val: clientID}}
}

// Streams hold the relevant information from the top-level JSON
type Streams struct {
	Error   string   `json:"error"`  // If json key "error" exists something must've gone wrong
	Total   float64  `json:"_total"` // json.Unmarshal need underscore to be in quotations inside quotations...
	Streams []Stream `json:"streams"`
}

// stream holds the relevant information in the "streams" key of the original JSON
type Stream struct {
	Id              float64           `json:"_id"`
	Viewers         float64           `json:"viewers"`
	PreviewImageUrl map[string]string `json:"preview"`
	Channel         ChannelInfo       `json:"channel"`
}

// channelInfo holds the relevant information in the "channel" key of the original JSON
type ChannelInfo struct {
	Id           float64 `json:"_id"`
	Name         string  `json:"name"`
	DisplayName  string  `json:"display_name"`
	LogoImageUrl string  `json:"logo"`
	Status       string  `json:"status"`
}

// twitch.Streams gets the top "limit" of "game" from the twitch api.
func (t Twitch) Streams(game, limit string) api.Streams {
	streams := Streams{}
	gameQ := api.Pair{"game", game}
	limitQ := api.Pair{"limit", limit}
	success := api.ApiCall(&streams, "https://api.twitch.tv/kraken/streams/", &t.auth, gameQ, limitQ)

	if success {
		return convertStreamInfo(&streams)
	} else {
		return api.Streams{
			Error:        true,
			ErrorMessage: "Failed to retreive information from API",
		}
	}
}

func convertStreamInfo(twitchStreams *Streams) (apiStreams api.Streams) {
	apiStreams.Total = twitchStreams.Total
	if twitchStreams.Error != "" { // Twitch errors is a json key "error" that holds the error message. This gets unmarshaled into "" or the message.
		apiStreams.Error = true
	} else {
		apiStreams.Error = false
	}
	apiStreams.ErrorMessage = twitchStreams.Error

	// Move creation of api.Stream to function?
	for _, v := range twitchStreams.Streams {
		apiStreams.Streams = append(apiStreams.Streams, api.Stream{
			Id:              v.Id,
			PreviewImageUrl: v.PreviewImageUrl["large"], // Twitch has a number of preview urls, one of which is "large"
			LogoImageUrl:    v.Channel.LogoImageUrl,
			Viewers:         v.Viewers,
			ChannelId:       v.Channel.Id,
			UserName:        v.Channel.Name,
			DisplayName:     v.Channel.DisplayName,
			StatusText:      v.Channel.Status})
	}
	return
}
