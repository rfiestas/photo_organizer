package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	photoorganizer "github.com/rfiestas/photo_organizer/internal/photo_organizer"
)

func printHelp() {
	fmt.Println(`Usage:
  photo_organizer --mode generate --folder <path>
    Generate .xmp metadata files for each image in the folder.

  photo_organizer --mode search --folder <path> --query <text>
    Search for images matching the query in the metadata.

Options:
  --mode     Operation mode: generate or search.
  --folder   Path to the folder containing image files.
  --query    Text to search in metadata (only used in search mode).
  -h         Show this help message.`)
}

func main() {
	mode := flag.String("mode", "generate", "Operation mode: generate or search.")
	folder := flag.String("folder", "", "Path to the folder containing images.")
	query := flag.String("query", "", "Search query (only used in search mode).")
	flag.Parse()

	if *mode == "" {
		printHelp()
		os.Exit(1)
	}

	if *folder == "" {
		fmt.Println("Error: --folder is required.")
		os.Exit(1)
	}

	// Use current path when folder input is not absolute.
	if !filepath.IsAbs(*folder) {
		path, err := os.Getwd()
		if err != nil {
			log.Println(err)
		}
		completePath := filepath.Join(path, *folder)
		folder = &completePath
	}

	switch *mode {
	case "generate":
		fmt.Println("Generating metadata for folder:", *folder)
		describer := photoorganizer.LlavaDescriber{}
		imageList, err := photoorganizer.FindImagesInFolder(*folder)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		photoorganizer.Generate(imageList, describer)
	case "search":
		if *query == "" {
			fmt.Println("Error: --query is required in search mode.")
			os.Exit(1)
		}
		// TODO: implement search logic
		fmt.Println("Search mode is not yet implemented.")
	default:
		fmt.Println("Error: invalid --mode value:", *mode)
		printHelp()
		os.Exit(1)
	}
}
