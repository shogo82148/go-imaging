package main

import (
	"image"
	"log"
	"os"

	"github.com/shogo82148/go-imaging/fp16"
	"github.com/shogo82148/go-imaging/icc"
	"github.com/shogo82148/go-imaging/png"
	"github.com/shogo82148/go-imaging/srgb"
)

func main() {
	if len(os.Args) != 4 {
		log.Println("usage: decodeTone input-png output-png icc-profile")
	}

	input := os.Args[1]
	output := os.Args[2]
	profile := os.Args[3]

	img, err := decodeTone(input)
	if err != nil {
		log.Fatal(err)
	}

	p, err := decodeProfile(profile)
	if err != nil {
		log.Fatal(err)
	}

	w, err := os.Create(output)
	if err != nil {
		log.Fatal(err)
	}
	err = png.EncodeWithMeta(w, &png.ImageWithMeta{
		Image:          img,
		ICCProfileName: "GIMP built-in Linear sRGB",
		ICCProfile:     p,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func decodeTone(filename string) (*fp16.NRGBAh, error) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	input, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	return srgb.DecodeTone(input), nil
}

func decodeProfile(filename string) (*icc.Profile, error) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	return icc.Decode(f)
}
