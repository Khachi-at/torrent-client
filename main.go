package main

import (
	"log"
	"os"

	"github.com/khachi-at/torrent-client/internal/torrentfile"
)

func main() {
	inPath := os.Args[1]
	outPath := os.Args[2]

	tf, err := torrentfile.Open(inPath)
	if err != nil {
		log.Fatal(err)
	}
	if err := tf.DownloadToFile(outPath); err != nil {
		log.Fatal(err)
	}
}
