package main

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"

	"github.com/EdlinOrg/prominentcolor"
	"github.com/gerow/go-color"
)

type HSB struct {
	H int64
	S int64
	B int64
}

type DetectedColours struct {
	Primary   *HSB
	Secondary *HSB
	Tertiary  *HSB
}

func getAlbumArtworkDominantColours(url string) (DetectedColours, error) {
	var detectedColours DetectedColours
	img, err := getImageFromUrl(url)
	if err != nil {
		return detectedColours, err
	}

	colours, err := prominentcolor.KmeansWithArgs(prominentcolor.ArgumentNoCropping, img)
	if err != nil {
		log.Fatalf("Failed to get prominent colours from album artwork. Error: %s", err)
	}

	for i, c := range colours {
		if i == 0 {
			detectedColours.Primary = convertHexToHSBK(c.AsString())
		}
		if i == 1 {
			detectedColours.Secondary = convertHexToHSBK(c.AsString())
		}
		if i == 2 {
			detectedColours.Tertiary = convertHexToHSBK(c.AsString())
		}
		if i >= 3 {
			break
		}
	}

	return detectedColours, err
}

func getImageFromUrl(url string) (image.Image, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download image data from URL %s. \nError: %s", url, err)
	}

	imgBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read data from URL response. \nError: %s", err)
	}

	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to decode raw image data into Image. \nError: %s", err)
	}

	return img, nil
}

func loadImage(fileInput string) (image.Image, error) {
	f, err := os.Open(fileInput)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}

func convertHexToHSBK(hex string) *HSB {
	rgb, _ := color.HTMLToRGB(hex)
	hsl := rgb.ToHSL()
	hue := int64(math.Round(65535*hsl.H)) % 65535
	sat := int64(math.Round(65535 * hsl.S))
	brightness := int64(math.Round(65535 * hsl.L))

	return &HSB{
		H: hue,
		S: sat,
		B: brightness,
	}
}
