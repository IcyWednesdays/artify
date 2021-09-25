package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/zmb3/spotify"
)

var html = `
<br/>
`

var (
	endpoint string
	user     *spotify.PrivateUser
)

func main() {
	endpoint, isSet := os.LookupEnv("ARTIFY_ENDPOINT")
	if !isSet {
		endpoint = "localhost"
	}

	// Callback for authenticating the user
	http.HandleFunc("/callback", handleAuthCallback)

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
