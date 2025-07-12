// package main

// import (
// 	"log"
// 	"os"

// 	"gioui.org/app"
// )

// type ProjectOpenedState struct {
// 	ProjectPath string
// }

// func main() {
// 	go func() {
// 		window := new(app.Window)
// 		err := run(window)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		os.Exit(0)
// 	}()
// 	app.Main()
// }

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"example.com/main/translator"
)

func main() {
	// Define command line flags
	dataPath := flag.String("data", "data", "Path to the data directory containing source files")
	dictionaryPath := flag.String("dictionary", "dictionary", "Path to the dictionary directory containing translation files")
	exportPath := flag.String("export", "export", "Path to the export directory for translated files")

	flag.Parse()

	// Validate paths exist
	if _, err := os.Stat(*dataPath); os.IsNotExist(err) {
		log.Fatalf("Data directory does not exist: %s", *dataPath)
	}

	if _, err := os.Stat(*dictionaryPath); os.IsNotExist(err) {
		log.Fatalf("Dictionary directory does not exist: %s", *dictionaryPath)
	}

	// Get the directory of the executable to resolve relative paths
	execDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}

	// Resolve relative paths
	if !filepath.IsAbs(*dataPath) {
		*dataPath = filepath.Join(execDir, *dataPath)
	}
	if !filepath.IsAbs(*dictionaryPath) {
		*dictionaryPath = filepath.Join(execDir, *dictionaryPath)
	}
	if !filepath.IsAbs(*exportPath) {
		*exportPath = filepath.Join(execDir, *exportPath)
	}

	// Create translator and run translation
	translatorInstance := translator.NewTranslator(*dataPath, *dictionaryPath, *exportPath)

	fmt.Printf("Starting translation process...\n")
	fmt.Printf("Data path: %s\n", *dataPath)
	fmt.Printf("Dictionary path: %s\n", *dictionaryPath)
	fmt.Printf("Export path: %s\n", *exportPath)

	err = translatorInstance.Translate()
	if err != nil {
		log.Fatalf("Translation failed: %v", err)
	}

	fmt.Printf("Translation completed successfully!\n")
	fmt.Printf("Translated files written to: %s\n", *exportPath)
}
