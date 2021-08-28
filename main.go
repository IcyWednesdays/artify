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
	user *spotify.PrivateUser
)

func main() {
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

	http.ListenAndServe(":8080", nil)

}
