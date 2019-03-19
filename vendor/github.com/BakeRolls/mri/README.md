# MangaRockImage

[![GoDoc](https://godoc.org/github.com/bakerolls/mri?status.svg)](https://godoc.org/github.com/bakerolls/mri)
[![Go Report Card](https://goreportcard.com/badge/github.com/bakerolls/mri)](https://goreportcard.com/report/github.com/bakerolls/mri)

Decode .mri files you might have found on MangaRock. As with other image formats, just import the package.

```bash
$ # Convert an image to PNG:
$ mri image.mri image.png
```

```go
package main

import (
	"image"
	"log"

	_ "github.com/bakerolls/mri"
)

func main() {
	// ...

	img, _, err := image.Decode(r)
	if err != nil {
		log.Fatalf("could not decode .mri: %v", err)
	}

	// ...
}
```
