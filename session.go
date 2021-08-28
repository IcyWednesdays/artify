package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/zmb3/spotify"
)

type ActiveSession struct {
	currentSongChannel      chan *CurrentSong
	client                  *spotify.Client
	sessionCompletedChannel chan bool
}

type CurrentSong struct {
	AlbumArtworkUrl string
	DetectedColours
	Name string
}

type NotRelevantReason int

const (
	IncorrectDeviceId NotRelevantReason = iota
	NoActiveSession
	Null
)

func isRelevantListeningEvent(playerState *spotify.PlayerState) (bool, NotRelevantReason) {
	if playerState.CurrentlyPlaying.Playing == true && playerState.Device.Name == os.Getenv("ARTIFY_PLAYER_DEVICE") {
		return true, Null
	}

	if playerState.CurrentlyPlaying.Playing == false {
		return false, NoActiveSession
	}

	return false, IncorrectDeviceId
}

func (a *ActiveSession) emitCurrentSong(currentSong *spotify.FullTrack) {
	albumArtColours, err := getAlbumArtworkDominantColours(currentSong.Album.Images[0].URL)
	if err != nil {
		log.Fatalf("Unable to get dominant colours of image at %s", currentSong.Album.Images[0].URL)
	}

	// Emit an ActiveSession obj to the channel containing the new song details
	a.currentSongChannel <- &CurrentSong{
		AlbumArtworkUrl: currentSong.Album.Images[0].URL,
		DetectedColours: albumArtColours,
		Name:            currentSong.Name,
	}
}

// Monitors a listening session and emits a CurrentSong via the currentSongChannel
func (a *ActiveSession) observeActiveSession(playerState *spotify.PlayerState) {
	playerCooldown := 0
	isCurrentlyPlaying := playerState.Playing

	a.emitCurrentSong(playerState.Item)
	currentSongId := playerState.CurrentlyPlaying.Item.ID

	// Continue while a session is continuing, or until 30 seconds after it finishes
	for maxPlayerCooldown := 7; playerCooldown < maxPlayerCooldown && isCurrentlyPlaying; {
		// Increment the counter if songs have stopped playing
		if !isCurrentlyPlaying {
			playerCooldown++
		} else {
			playerCooldown = 0
		}

		time.Sleep(5 * time.Second)

		// Get the current state of the player
		currentPlayerState, err := a.client.PlayerState()
		if err != nil {
			log.Fatalf("Failed to get current player state. Error: %s", err)
		}

		// If the now playing song is different from the previous
		if currentPlayerState.Item.ID != currentSongId {
			// Emit the current song to the activeSessionChannel
			a.emitCurrentSong(currentPlayerState.Item)
			// Update currentSongId with the ID of the current song
			currentSongId = currentPlayerState.Item.ID
		}

		isCurrentlyPlaying, _ = isRelevantListeningEvent(currentPlayerState)
	}

	a.sessionCompletedChannel <- true
	close(a.currentSongChannel)
}

// Used so we can log errors out directly rather than handling them in the
func getPlayerState(client *spotify.Client) *spotify.PlayerState {
	playerState, err := client.PlayerState()
	if err != nil {
		log.Fatalf("Unable to get player state for user %s. Error: %s", user.ID, err)
	}

	return playerState
}

func monitorForActiveSession(client *spotify.Client) {
	for {
		playerState, err := client.PlayerState()
		if err != nil {
			log.Fatalf("Unable to get player state for user %s. Error: %s", user.ID, err)
		}
		// Check for relevant sessions every 30 seconds, and break once one is found
		for isRelevant, _ := isRelevantListeningEvent(playerState); !isRelevant; isRelevant, _ = isRelevantListeningEvent(getPlayerState(client)) {
			fmt.Println("No active session found. Waiting...")
			time.Sleep(10 * time.Second) // TODO: Check less often overnight/low usage times
		}

		activeSession := &ActiveSession{
			currentSongChannel:      make(chan *CurrentSong),
			client:                  client,
			sessionCompletedChannel: make(chan bool, 1),
		}

		// Execute the observe routine, and the routine that updates the lights. These routines communicate via the currentSongChannel
		go activeSession.observeActiveSession(playerState)
		go activeSession.updateLightColours()

		// Once the routines are started, wait for a completion message from the sessionCompletedChannel. Once this is received, break so the whole monitoring loop restarts
		for activeSessionComplete := range activeSession.sessionCompletedChannel {
			if activeSessionComplete {
				break
			}
		}
	}
}
