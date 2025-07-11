package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

// ParsedReference stores parsed components
type ParsedReference struct {
	Tag        string
	ItemName   string
	SourceBook string
	FluffText  string
}

func parseReferencesFromJSON(filePath string) ([]ParsedReference, error) {
	var results []ParsedReference

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var jsonContent interface{}
	if err := json.Unmarshal(data, &jsonContent); err != nil {
		return nil, err
	}

	// Case 1: Tag-based pattern: {@tag name|source|fluff}
	taggedPattern := regexp.MustCompile(`{@(\w+)\s([^|{}]+)(?:\|([^|{}]+))?(?:\|([^{}]+))?}`)

	// Case 2: Simple pattern: item|source
	simplePattern := regexp.MustCompile(`^([^|{}]+)\|([^|{}]+)$`)

	// Recursive walker
	var walk func(interface{})
	walk = func(node interface{}) {
		switch val := node.(type) {
		case map[string]interface{}:
			for _, v := range val {
				walk(v)
			}
		case []interface{}:
			for _, item := range val {
				walk(item)
			}
		case string:
			// Try tagged pattern
			for _, match := range taggedPattern.FindAllStringSubmatch(val, -1) {
				if len(match) >= 3 {
					ref := ParsedReference{
						Tag:        match[1],
						ItemName:   match[2],
						SourceBook: "",
						FluffText:  "",
					}
					if len(match) >= 4 {
						ref.SourceBook = match[3]
					}
					if len(match) == 5 {
						ref.FluffText = match[4]
					}
					if ref.Tag == "b" || ref.Tag == "book" {
						continue
					}
					results = append(results, ref)
				}
			}
			// Try simple pattern (only if not inside tagged form)
			if match := simplePattern.FindStringSubmatch(val); match != nil {
				ref := ParsedReference{
					Tag:        "",
					ItemName:   match[1],
					SourceBook: match[2],
					FluffText:  "",
				}
				results = append(results, ref)
			}
		}
	}

	walk(jsonContent)
	return results, nil
}

func main() {
	refs, err := parseReferencesFromJSON("C:\\Users\\Maksym\\Desktop\\plutonium\\data\\backgrounds.json")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	for _, ref := range refs {
		fmt.Printf("Tag: %s | Name: %s | Source: %s | Fluff: %s\n",
			ref.Tag, ref.ItemName, ref.SourceBook, ref.FluffText)
	}
}
