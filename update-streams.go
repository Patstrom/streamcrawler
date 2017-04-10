package main

import (
	"flag"
	"fmt"
	"strings"
    "database/sql"
	"repos/patstrom/update-streams/api"
	"repos/patstrom/update-streams/hitbox"
	"repos/patstrom/update-streams/twitch"
)

var site, env string
var connection api.API
var gMap map[string]string // Maps generic name in "games" to what the api expects. Probably not needed if we get game_id from database
var database *sql.DB // Database connection
var stmtOnlineStreams, stmtUpdateStream, stmtStreamExists, stmtInsertStream *sql.Stmt // Queries

// Games to check
var games = []string{"overwatch", "dota", "cs", "lol"}

func init() {
	const (
		defaultSite = "twitch"
		defaultEnv  = "test"
		siteUsage   = "The site which API to target"
		envUsage    = "The database database environment to target"
	)
	flag.StringVar(&site, "site", defaultSite, siteUsage)
	flag.StringVar(&site, "s", defaultSite, siteUsage+" (shorthand)")

	flag.StringVar(&env, "env", defaultEnv, envUsage)
	flag.StringVar(&env, "e", defaultEnv, envUsage+" (shorthand)")
}

// Setup api by creating a new api.API object corresponding to the api we are currently targeting
func init() {
	flag.Parse()
	switch {
	case strings.EqualFold("twitch", site):
		connection = twitch.NewTwitch("ppq86qhn55lqxfotpqcz69phndjmxg") // Should this be read from file?
		// A map for the generic game names to what the API expects.
		// Should this be defines in the corresponding package? E.g twitch.Map(). Probably not
		gMap = map[string]string{
			"overwatch": "Overwatch",
			"dota":      "Dota 2",
			"cs":        "Counter-Strike: Global Offensive",
			"lol":       "League of Legends",
		}
	case strings.EqualFold("hitbox", site):
		connection = &hitbox.Hitbox{}
		gMap = map[string]string{
			"overwatch": "overwatch",
			"dota":      "dota-2",
			"cs":        "counter-strike-global-offensive",
			"lol":       "league-of-legends",
		}
	default:
        // Log.fatal this
		panic("API not implemented")
	}
}

/*
// setup database by creating a new sql.DB object corresponding to the environment we are currently targeting
func init() {
    database, err := sql.Open("") // Open database, this should correspond to the environment I suppose
    if err != nil {
        // Log.fatal this
    }

    stmtOnlineStreams, err := database.Prepare("") // Get already online streams and parse them into api.Streams
    if err != nil {
        // Log.fatal this
    }

    stmtUpdateStream, err := database.Prepare("") // Update stream with new info
    if err != nil {
        // Log.fatal this
    }

    stmtStreamExists, err := database.Prepare("") // See if stream exists in database (to see whether we need to update online status or insert a new stream)
    if err != nil {
        // Log.fatal this
    }

    stmtInsertStream, err := database.Prepare("") // Insert new stream
    if err != nil {
        // Log.fatal this
    }
}
*/

func main() {
	for _, game := range games {
		// Get all we have on "game" in db
		streams := connection.Streams(gMap[game], "1") // Get top 100 of "game"

		fmt.Println("Game:", gMap[game], "Total:", streams.Total, "Error:", streams.Error, "Error Message:", streams.ErrorMessage)
		for _, v := range streams.Streams {
			fmt.Printf("Id: %v\nViewers: %v\nPreview: %v\nLogo: %v\nChannel: {\n\tChannelId: %v\n\tUsername: %v\n\tDisplayName: %v\n\tStatusText: %v\n}\n", v.Id, v.Viewers, v.PreviewImageUrl, v.LogoImageUrl, v.ChannelId, v.UserName, v.DisplayName, v.StatusText)
		}

        /*
            if v in dbStreams
                if metadata is updated // Is this even worth it? Probably not, due to preview and stuff. Just update
                    continue
                else
                    stmtUpdateStream
                remove the object in dbStreams // This is important
            else
                if stmtStreamExists
                    stmtUpdateStream
                else
                    stmtInsertStream

            for range dbStreams
                stmtUpdateStreams to offline
        */

	}

    /*
    // Close stuff
    stmtOnlineStreams.Close()
    stmtUpdateStream.Close()
    stmtStreamExists.Close()
    stmtInsertStream.Close()
    database.Close()
    */
}
