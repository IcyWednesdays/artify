package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zmb3/spotify"
)

var html = `
<br/>
`

var (
	authUrlChan = make(chan string, 1)
	endpoint    = getEnv("ARTIFY_ENDPOINT", "localhost")
	user        *spotify.PrivateUser
)

func main() {
	// Callback for authenticating the user
	http.HandleFunc("/callback", handleAuthCallback)

	// Receive the auth URL and create a redirect on /auth
	http.HandleFunc("/auth", authRedirectHandler)

	// Once authenticated, kick off the main goroutine for the program
	go func() {
		client := getClientForUser()
		user, err := client.CurrentUser()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Logged in as:", user.ID)

		go monitorForActiveSession(client)
	}()

	http.ListenAndServe(fmt.Sprintf("%s:8080", endpoint), nil)
}
