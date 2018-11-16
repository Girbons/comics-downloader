package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func usage() {
	fmt.Println("comics-download -u=http://comic-source")
	os.Exit(0)
}

func main() {
	url := flag.String("url", "", "Comic URL")

	flag.Usage = usage
	flag.Parse()

	if *url == "" {
		log.Fatal("url parameter is required")
	}

	DetectComic(*url)
	fmt.Println("Download Completed")
}
