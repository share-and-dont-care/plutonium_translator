package translator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Translator handles the translation process
type Translator struct {
	dataPath       string
	dictionaryPath string
	exportPath     string
}

// NewTranslator creates a new translator instance
func NewTranslator(dataPath, dictionaryPath, exportPath string) *Translator {
	return &Translator{
		dataPath:       dataPath,
		dictionaryPath: dictionaryPath,
		exportPath:     exportPath,
	}
}

// Translate processes all background files and applies translations
func (t *Translator) Translate() error {
	// Read dictionary data
	dictionaryEntries, err := t.loadDictionaryEntries()
	if err != nil {
		return fmt.Errorf("failed to load dictionary data: %w", err)
	}

	// Read source data
	sourceData, err := t.loadSourceData()
	if err != nil {
		return fmt.Errorf("failed to load source data: %w", err)
	}

	// Apply translations
	translatedData, err := t.applyTranslations(sourceData, dictionaryEntries)
	if err != nil {
		return fmt.Errorf("failed to apply translations: %w", err)
	}

	// Write translated data to export directory
	err = t.writeTranslatedData(translatedData)
	if err != nil {
		return fmt.Errorf("failed to write translated data: %w", err)
	}

	return nil
}

// loadDictionaryEntries loads all dictionary files and returns a slice of entries
func (t *Translator) loadDictionaryEntries() ([]map[string]interface{}, error) {
	files, err := filepath.Glob(filepath.Join(t.dictionaryPath, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to glob dictionary files: %w", err)
	}

	var allEntries []map[string]interface{}

	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read dictionary file %s: %w", file, err)
		}

		var dictData map[string]interface{}
		err = json.Unmarshal(data, &dictData)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal dictionary file %s: %w", file, err)
		}

		// Extract background entries from the dictionary
		if backgroundEntry, exists := dictData["background"]; exists {
			switch v := backgroundEntry.(type) {
			case map[string]interface{}:
				// Single object format
				allEntries = append(allEntries, v)
			case []interface{}:
				// Array format
				for _, item := range v {
					if backgroundMap, ok := item.(map[string]interface{}); ok {
						allEntries = append(allEntries, backgroundMap)
					}
				}
			}
		}
	}

	return allEntries, nil
}

// loadSourceData loads the source background data
func (t *Translator) loadSourceData() (map[string]interface{}, error) {
	data, err := ioutil.ReadFile(filepath.Join(t.dataPath, "backgrounds.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to read source data: %w", err)
	}

	var sourceData map[string]interface{}
	err = json.Unmarshal(data, &sourceData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal source data: %w", err)
	}

	return sourceData, nil
}

// applyTranslations applies dictionary translations to source data
func (t *Translator) applyTranslations(sourceData map[string]interface{}, dictionaryEntries []map[string]interface{}) (map[string]interface{}, error) {
	// Create a copy of source data
	translatedData := make(map[string]interface{})
	for k, v := range sourceData {
		translatedData[k] = v
	}

	// Get the backgrounds array from source data
	backgroundsInterface, exists := sourceData["background"]
	if !exists {
		return translatedData, fmt.Errorf("no 'background' field found in source data")
	}

	backgroundsArray, ok := backgroundsInterface.([]interface{})
	if !ok {
		return translatedData, fmt.Errorf("'background' field is not an array")
	}

	fmt.Printf("Found %d source backgrounds\n", len(backgroundsArray))
	fmt.Printf("Found %d dictionary entries\n", len(dictionaryEntries))

	// Create a map for quick lookup of source backgrounds
	sourceBackgroundsMap := make(map[string]interface{})
	for _, bg := range backgroundsArray {
		if bgMap, ok := bg.(map[string]interface{}); ok {
			name, nameOk := bgMap["name"].(string)
			source, sourceOk := bgMap["source"].(string)
			if nameOk && sourceOk {
				key := name + "|" + source
				sourceBackgroundsMap[key] = bgMap
				fmt.Printf("Source background: %s|%s\n", name, source)
			}
		}
	}

	// Process each dictionary entry
	var translatedBackgrounds []interface{}
	for _, dictEntry := range dictionaryEntries {
		originName, originNameOk := dictEntry["origin_name"].(string)
		originSource, originSourceOk := dictEntry["origin_source"].(string)

		if !originNameOk || !originSourceOk {
			fmt.Printf("Skipping dictionary entry without proper origin info\n")
			continue // Skip entries without proper origin info
		}

		fmt.Printf("Dictionary entry: %s|%s\n", originName, originSource)

		// Find matching source background
		key := originName + "|" + originSource
		if sourceBackground, exists := sourceBackgroundsMap[key]; exists {
			fmt.Printf("Found match for: %s\n", key)
			// Create a copy of the source background
			sourceBgMap := sourceBackground.(map[string]interface{})
			translatedBg := make(map[string]interface{})

			// Copy all properties from source
			for k, v := range sourceBgMap {
				translatedBg[k] = v
			}

			// Apply translations from dictionary
			if translatedName, exists := dictEntry["name"]; exists {
				translatedBg["name"] = translatedName
			}

			// Apply entries translation if present
			if dictEntries, exists := dictEntry["entries"]; exists {
				translatedBg["entries"] = dictEntries
			}

			translatedBackgrounds = append(translatedBackgrounds, translatedBg)
		} else {
			fmt.Printf("No match found for: %s\n", key)
		}
	}

	fmt.Printf("Created %d translated backgrounds\n", len(translatedBackgrounds))

	// Update the translated data with the processed backgrounds
	translatedData["background"] = translatedBackgrounds

	return translatedData, nil
}

// writeTranslatedData writes the translated data to the export directory
func (t *Translator) writeTranslatedData(data map[string]interface{}) error {
	// Ensure export directory exists
	err := os.MkdirAll(t.exportPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create export directory: %w", err)
	}

	// Marshal the data
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal translated data: %w", err)
	}

	// Write to file
	outputPath := filepath.Join(t.exportPath, "backgrounds.json")
	err = ioutil.WriteFile(outputPath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write translated data: %w", err)
	}

	return nil
}
