package translator

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNewTranslator(t *testing.T) {
	translator := NewTranslator("data", "dictionary", "export")

	if translator.dataPath != "data" {
		t.Errorf("Expected dataPath to be 'data', got '%s'", translator.dataPath)
	}

	if translator.dictionaryPath != "dictionary" {
		t.Errorf("Expected dictionaryPath to be 'dictionary', got '%s'", translator.dictionaryPath)
	}

	if translator.exportPath != "export" {
		t.Errorf("Expected exportPath to be 'export', got '%s'", translator.exportPath)
	}
}

func TestLoadDictionaryEntries(t *testing.T) {
	// Create temporary test directory
	tempDir, err := ioutil.TempDir("", "test_dict")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test dictionary file
	dictData := map[string]interface{}{
		"background": map[string]interface{}{
			"origin_name":   "Acolyte",
			"origin_source": "XPHB",
			"name":          "Аколіт [Acolyte]",
			"entries": []interface{}{
				map[string]interface{}{
					"items": []interface{}{
						map[string]interface{}{
							"name":  "Здібності:",
							"entry": "Інтелект, Мудрість, Харизма",
						},
						map[string]interface{}{
							"name":  "Риси:",
							"entry": "{@feat Magic Initiate|XPHB} (Клірик)",
						},
					},
				},
			},
		},
	}

	dictJSON, err := json.MarshalIndent(dictData, "", "    ")
	if err != nil {
		t.Fatalf("Failed to marshal test dictionary data: %v", err)
	}

	dictPath := filepath.Join(tempDir, "backgrounds.json")
	err = ioutil.WriteFile(dictPath, dictJSON, 0644)
	if err != nil {
		t.Fatalf("Failed to write test dictionary file: %v", err)
	}

	translator := NewTranslator("", tempDir, "")
	dictionaryEntries, err := translator.loadDictionaryEntries()

	if err != nil {
		t.Fatalf("Failed to load dictionary data: %v", err)
	}

	if len(dictionaryEntries) != 1 {
		t.Errorf("Expected 1 dictionary entry, got %d", len(dictionaryEntries))
	}

	entry := dictionaryEntries[0]
	if entry["origin_name"] != "Acolyte" {
		t.Errorf("Expected origin_name to be 'Acolyte', got '%s'", entry["origin_name"])
	}

	if entry["origin_source"] != "XPHB" {
		t.Errorf("Expected origin_source to be 'XPHB', got '%s'", entry["origin_source"])
	}

	if entry["name"] != "Аколіт [Acolyte]" {
		t.Errorf("Expected name to be 'Аколіт [Acolyte]', got '%s'", entry["name"])
	}
}

func TestLoadSourceData(t *testing.T) {
	// Create temporary test directory
	tempDir, err := ioutil.TempDir("", "test_data")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test source data
	sourceData := map[string]interface{}{
		"_meta": map[string]interface{}{
			"internalCopies": []interface{}{"background"},
		},
		"background": []interface{}{
			map[string]interface{}{
				"name":   "Acolyte",
				"source": "XPHB",
				"page":   178,
				"srd52":  true,
			},
		},
	}

	sourceJSON, err := json.MarshalIndent(sourceData, "", "    ")
	if err != nil {
		t.Fatalf("Failed to marshal test source data: %v", err)
	}

	sourcePath := filepath.Join(tempDir, "backgrounds.json")
	err = ioutil.WriteFile(sourcePath, sourceJSON, 0644)
	if err != nil {
		t.Fatalf("Failed to write test source file: %v", err)
	}

	translator := NewTranslator(tempDir, "", "")
	loadedData, err := translator.loadSourceData()

	if err != nil {
		t.Fatalf("Failed to load source data: %v", err)
	}

	backgrounds, exists := loadedData["background"]
	if !exists {
		t.Fatalf("Expected 'background' field in loaded data")
	}

	backgroundsArray, ok := backgrounds.([]interface{})
	if !ok {
		t.Fatalf("Expected 'background' to be an array")
	}

	if len(backgroundsArray) != 1 {
		t.Errorf("Expected 1 background, got %d", len(backgroundsArray))
	}

	background := backgroundsArray[0].(map[string]interface{})
	if background["name"] != "Acolyte" {
		t.Errorf("Expected background name to be 'Acolyte', got '%s'", background["name"])
	}

	if background["source"] != "XPHB" {
		t.Errorf("Expected background source to be 'XPHB', got '%s'", background["source"])
	}
}

func TestApplyTranslations(t *testing.T) {
	// Create source data
	sourceData := map[string]interface{}{
		"_meta": map[string]interface{}{
			"internalCopies": []interface{}{"background"},
		},
		"background": []interface{}{
			map[string]interface{}{
				"name":   "Acolyte",
				"source": "XPHB",
				"page":   178,
				"srd52":  true,
			},
		},
	}

	// Create dictionary data
	dictionaryEntries := []map[string]interface{}{
		{
			"origin_name":   "Acolyte",
			"origin_source": "XPHB",
			"name":          "Аколіт [Acolyte]",
			"entries": []interface{}{
				map[string]interface{}{
					"items": []interface{}{
						map[string]interface{}{
							"name":  "Здібності:",
							"entry": "Інтелект, Мудрість, Харизма",
						},
					},
				},
			},
		},
	}

	translator := NewTranslator("", "", "")
	translatedData, err := translator.applyTranslations(sourceData, dictionaryEntries)

	if err != nil {
		t.Fatalf("Failed to apply translations: %v", err)
	}

	backgrounds, exists := translatedData["background"]
	if !exists {
		t.Fatalf("Expected 'background' field in translated data")
	}

	backgroundsArray, ok := backgrounds.([]interface{})
	if !ok {
		t.Fatalf("Expected 'background' to be an array")
	}

	if len(backgroundsArray) != 1 {
		t.Errorf("Expected 1 background, got %d", len(backgroundsArray))
	}

	translatedBackground := backgroundsArray[0].(map[string]interface{})
	if translatedBackground["name"] != "Аколіт [Acolyte]" {
		t.Errorf("Expected translated name to be 'Аколіт [Acolyte]', got '%s'", translatedBackground["name"])
	}

	// Check that other fields remain unchanged
	if translatedBackground["source"] != "XPHB" {
		t.Errorf("Expected source to remain 'XPHB', got '%s'", translatedBackground["source"])
	}

	if translatedBackground["page"] != float64(178) {
		t.Errorf("Expected page to remain 178, got %v", translatedBackground["page"])
	}

	if translatedBackground["srd52"] != true {
		t.Errorf("Expected srd52 to remain true")
	}
}

func TestWriteTranslatedData(t *testing.T) {
	// Create temporary test directory
	tempDir, err := ioutil.TempDir("", "test_export")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test data
	data := map[string]interface{}{
		"_meta": map[string]interface{}{
			"internalCopies": []interface{}{"background"},
		},
		"background": []interface{}{
			map[string]interface{}{
				"name":   "Аколіт [Acolyte]",
				"source": "XPHB",
				"page":   178,
				"srd52":  true,
			},
		},
	}

	translator := NewTranslator("", "", tempDir)
	err = translator.writeTranslatedData(data)

	if err != nil {
		t.Fatalf("Failed to write translated data: %v", err)
	}

	// Verify file was created
	outputPath := filepath.Join(tempDir, "backgrounds.json")
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Expected output file to exist at %s", outputPath)
	}

	// Read and verify content
	content, err := ioutil.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	var readData map[string]interface{}
	err = json.Unmarshal(content, &readData)
	if err != nil {
		t.Fatalf("Failed to unmarshal output file: %v", err)
	}

	backgrounds, exists := readData["background"]
	if !exists {
		t.Fatalf("Expected 'background' field in output")
	}

	backgroundsArray, ok := backgrounds.([]interface{})
	if !ok {
		t.Fatalf("Expected 'background' to be an array")
	}

	if len(backgroundsArray) != 1 {
		t.Errorf("Expected 1 background in output, got %d", len(backgroundsArray))
	}

	background := backgroundsArray[0].(map[string]interface{})
	if background["name"] != "Аколіт [Acolyte]" {
		t.Errorf("Expected background name to be 'Аколіт [Acolyte]', got '%s'", background["name"])
	}
}

func TestTranslateIntegration(t *testing.T) {
	// Create temporary test directories
	tempDataDir, err := ioutil.TempDir("", "test_data")
	if err != nil {
		t.Fatalf("Failed to create temp data directory: %v", err)
	}
	defer os.RemoveAll(tempDataDir)

	tempDictDir, err := ioutil.TempDir("", "test_dict")
	if err != nil {
		t.Fatalf("Failed to create temp dictionary directory: %v", err)
	}
	defer os.RemoveAll(tempDictDir)

	tempExportDir, err := ioutil.TempDir("", "test_export")
	if err != nil {
		t.Fatalf("Failed to create temp export directory: %v", err)
	}
	defer os.RemoveAll(tempExportDir)

	// Create test source data
	sourceData := map[string]interface{}{
		"_meta": map[string]interface{}{
			"internalCopies": []interface{}{"background"},
		},
		"background": []interface{}{
			map[string]interface{}{
				"name":   "Acolyte",
				"source": "XPHB",
				"page":   178,
				"srd52":  true,
			},
		},
	}

	sourceJSON, err := json.MarshalIndent(sourceData, "", "    ")
	if err != nil {
		t.Fatalf("Failed to marshal test source data: %v", err)
	}

	sourcePath := filepath.Join(tempDataDir, "backgrounds.json")
	err = ioutil.WriteFile(sourcePath, sourceJSON, 0644)
	if err != nil {
		t.Fatalf("Failed to write test source file: %v", err)
	}

	// Create test dictionary data
	dictData := map[string]interface{}{
		"background": map[string]interface{}{
			"origin_name":   "Acolyte",
			"origin_source": "XPHB",
			"name":          "Аколіт [Acolyte]",
			"entries": []interface{}{
				map[string]interface{}{
					"items": []interface{}{
						map[string]interface{}{
							"name":  "Здібності:",
							"entry": "Інтелект, Мудрість, Харизма",
						},
					},
				},
			},
		},
	}

	dictJSON, err := json.MarshalIndent(dictData, "", "    ")
	if err != nil {
		t.Fatalf("Failed to marshal test dictionary data: %v", err)
	}

	dictPath := filepath.Join(tempDictDir, "backgrounds.json")
	err = ioutil.WriteFile(dictPath, dictJSON, 0644)
	if err != nil {
		t.Fatalf("Failed to write test dictionary file: %v", err)
	}

	// Run translation
	translator := NewTranslator(tempDataDir, tempDictDir, tempExportDir)
	err = translator.Translate()

	if err != nil {
		t.Fatalf("Failed to translate: %v", err)
	}

	// Verify output
	outputPath := filepath.Join(tempExportDir, "backgrounds.json")
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Expected output file to exist at %s", outputPath)
	}

	content, err := ioutil.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	var outputData map[string]interface{}
	err = json.Unmarshal(content, &outputData)
	if err != nil {
		t.Fatalf("Failed to unmarshal output file: %v", err)
	}

	backgrounds, exists := outputData["background"]
	if !exists {
		t.Fatalf("Expected 'background' field in output")
	}

	backgroundsArray, ok := backgrounds.([]interface{})
	if !ok {
		t.Fatalf("Expected 'background' to be an array")
	}

	if len(backgroundsArray) != 1 {
		t.Errorf("Expected 1 background in output, got %d", len(backgroundsArray))
	}

	background := backgroundsArray[0].(map[string]interface{})
	if background["name"] != "Аколіт [Acolyte]" {
		t.Errorf("Expected background name to be 'Аколіт [Acolyte]', got '%s'", background["name"])
	}
}
