# MangaRockImage

[![GoDoc](https://godoc.org/github.com/bake/mri?status.svg)](https://godoc.org/github.com/bake/mri)
[![Go Report Card](https://goreportcard.com/badge/github.com/bake/mri)](https://goreportcard.com/report/github.com/bake/mri)

Decode .mri files you might have found on MangaRock. As with other image formats, just import the package.

```bash
$ # Convert an image to PNG:
$ mri image.mri image.png
$
```

```go
package main

import (
  "image"
  "log"

  _ "github.com/bake/mri"
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
