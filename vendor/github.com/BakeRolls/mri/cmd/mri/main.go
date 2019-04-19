package main

import (
	"image"
	"image/png"
	"log"
	"os"
	"strings"

	_ "github.com/bake/mri"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: mri image.mri [image.png]")
	}

	in := os.Args[1]
	r, err := os.Open(in)
	if err != nil {
		log.Fatalf("could not open file %s: %v", in, err)
	}
	img, _, err := image.Decode(r)
	if err != nil {
		log.Fatalf("could not decode .mri: %v", err)
	}

	out := strings.TrimSuffix(os.Args[1], ".mri") + ".png"
	if len(os.Args) >= 3 {
		out = os.Args[2]
	}
	w, err := os.Create(out)
	if err != nil {
		log.Fatalf("could not create file %s: %v", out, err)
	}
	if err := png.Encode(w, img); err != nil {
		log.Fatalf("could not encode png: %v", err)
	}
}
