package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/IcyWednesdays/golifx"
)

// Declare this as our own struct so we can add methods to it
type Bulb struct {
	*golifx.Bulb
}

// Listens to  activeSessionChannel and updates the lights whenever a change is received
func (a *ActiveSession) updateLightColours() {
	for currentSong := range a.currentSongChannel {
		fmt.Printf("Updating lights. Current song is %s\n", currentSong.Name)
		bulbs, err := golifx.LookupBulbs()
		if err != nil {
			log.Fatalf("Failed to look up bulbs. Error: %s", err)
		}

		// If we don't get 2 bulbs returned, try a few more times before giving up
		if len(bulbs) < 2 {
			attempts := 1
			maxAttempts := 3
			for bulbCount := len(bulbs); bulbCount < 2 && attempts < maxAttempts; attempts++ {
				if len(bulbs) == 2 {
					break
				}

				log.Printf("%d bulbs found. Retrying (%d/%d)", len(bulbs), attempts, maxAttempts)
				bulbs, err = golifx.LookupBulbs()

				time.Sleep(1 * time.Second)
			}
		}

		sceneConfiguration := getSceneByName(os.Getenv("ARTIFY_SCENE_NAME"))

		for _, b := range bulbs {
			bulb := Bulb{b}
			bulb.SetLightColour(currentSong.DetectedColours, sceneConfiguration)
		}
	}
}

func (b *Bulb) SetLightColour(detectedColours DetectedColours, sceneConfiguration *ZoneConfig) {
	for _, z := range sceneConfiguration.Primary {
		b.MultizoneSetColorZones(&golifx.HSBK{
			Hue:        uint16(detectedColours.Primary.H),
			Saturation: uint16(detectedColours.Primary.S),
			Brightness: uint16(detectedColours.Primary.B),
			Kelvin:     uint16(3500),
		}, uint32(5000), uint8(z.Start), uint8(z.End))
	}

	for _, z := range sceneConfiguration.Secondary {
		b.MultizoneSetColorZones(&golifx.HSBK{
			Hue:        uint16(detectedColours.Secondary.H),
			Saturation: uint16(detectedColours.Secondary.S),
			Brightness: uint16(detectedColours.Secondary.B),
			Kelvin:     uint16(3500),
		}, uint32(5000), uint8(z.Start), uint8(z.End))
	}

	for _, z := range sceneConfiguration.Tertiary {
		b.MultizoneSetColorZones(&golifx.HSBK{
			Hue:        uint16(detectedColours.Tertiary.H),
			Saturation: uint16(detectedColours.Tertiary.S),
			Brightness: uint16(detectedColours.Tertiary.B),
			Kelvin:     uint16(3500),
		}, uint32(5000), uint8(z.Start), uint8(z.End))
	}

}
