package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zmb3/spotify"
)

var (
	authChannel   = make(chan *spotify.Client, 1)
	authenticator = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadCurrentlyPlaying, spotify.ScopeUserReadPlaybackState, spotify.ScopeUserModifyPlaybackState)
	redirectURI   = "http://localhost:8080/callback"
	state         = "authState"
)

// Shows the auth prompt for the user, and returns the authenticated client back from the channel
func getClientForUser() *spotify.Client {
	defer close(authChannel)

	url := authenticator.AuthURL(state)
	fmt.Println("Log in to Spotify to continue:", url)

	return <-authChannel
}

func handleAuthCallback(w http.ResponseWriter, r *http.Request) {
	tok, err := authenticator.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := authenticator.NewClient(tok)
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "Login Completed!"+html)
	authChannel <- &client
}
