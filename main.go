package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strings"

	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

func main() {
	// Display the styled header of the program
	displayHeader()

	// Define the names of the YAML files to be compared
	file1Path := "examples/file1.yaml"
	file2Path := "examples/file2.yaml"

	// Read the content of the YAML files
	file1Content, err := readFile(file1Path)
	if err != nil {
		log.Fatalf("Error reading file %s: %v", file1Path, err)
	}

	file2Content, err := readFile(file2Path)
	if err != nil {
		log.Fatalf("Error reading file %s: %v", file2Path, err)
	}

	// Parse the content of the files into YAML documents
	yamlDocs1 := parseYAMLDocuments(file1Content)
	yamlDocs2 := parseYAMLDocuments(file2Content)

	// Create a mapping of the documents by "kind" and "metadata.name"
	yamlMap1 := mapYAMLDocuments(yamlDocs1)
	yamlMap2 := mapYAMLDocuments(yamlDocs2)

	// Compare the YAML documents and display the differences found
	fmt.Printf("\nDifferences found between the YAML files:\n\n")
	compareYAMLMaps(yamlMap1, yamlMap2)
}

// Display the styled header of the program
func displayHeader() {
	blue := color.New(color.FgBlue, color.Bold).SprintFunc()
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	fmt.Println(blue("========================================"))
	fmt.Println(green("             go-yaml-diff               "))
	fmt.Println(blue("========================================"))
}

// Read the content of a file and return it as a string
func readFile(filePath string) (string, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Parse multiple YAML documents into a slice of interfaces
func parseYAMLDocuments(content string) []interface{} {
	var documents []interface{}
	decoder := yaml.NewDecoder(strings.NewReader(content))
	for {
		var doc interface{}
		if err := decoder.Decode(&doc); err != nil {
			break
		}
		documents = append(documents, doc)
	}
	return documents
}

// Map YAML documents by "kind" and "metadata.name" to create a unique key for each object
func mapYAMLDocuments(documents []interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, doc := range documents {
		if yamlObject, ok := doc.(map[string]interface{}); ok {
			kind := yamlObject["kind"]
			metadata, hasMetadata := yamlObject["metadata"].(map[string]interface{})
			name := ""
			if hasMetadata {
				name = metadata["name"].(string)
			}
			if kind != nil && name != "" {
				key := fmt.Sprintf("%s/%s", kind, name)
				result[key] = yamlObject
			}
		}
	}
	return result
}

// Compare two maps of YAML documents and print differences with context
func compareYAMLMaps(yamlMap1, yamlMap2 map[string]interface{}) {
	for key, obj1 := range yamlMap1 {
		if obj2, exists := yamlMap2[key]; exists {
			fmt.Printf("Comparing object %s:\n", key)
			compareYAMLWithContext(obj1, obj2, "", key)
		} else {
			fmt.Printf("Object missing in file 2: %s\n", key)
		}
		fmt.Println("---")
	}

	// Check for objects that are only present in the second map
	for key := range yamlMap2 {
		if _, exists := yamlMap1[key]; !exists {
			fmt.Printf("Object missing in file 1: %s\n", key)
			fmt.Println("---")
		}
	}
}

// Compare YAML data and print differences with context, avoiding redundant information
func compareYAMLWithContext(data1, data2 interface{}, indent, parentPath string) {
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	map1, isMap1 := data1.(map[string]interface{})
	map2, isMap2 := data2.(map[string]interface{})

	// Compare maps (objects) and identify differences
	if isMap1 && isMap2 {
		for key := range map1 {
			value1 := map1[key]
			value2, exists := map2[key]

			currentPath := fmt.Sprintf("%s.%s", parentPath, key)
			if exists && !reflect.DeepEqual(value1, value2) {
				if isLeafNode(value1) && isLeafNode(value2) {
					fmt.Printf("Found difference on [%s]\n", currentPath)
					printDifferenceContext(red(fmt.Sprintf("- %s: %v", key, value1)), green(fmt.Sprintf("+ %s: %v", key, value2)), indent)
				} else {
					compareYAMLWithContext(value1, value2, indent+"  ", currentPath)
				}
			} else if !exists {
				fmt.Printf("%s%s\n", red(fmt.Sprintf("- %s: %v", key, formatYAMLValue(value1, indent))))
			}
		}
	}

	// Compare lists (arrays) and identify differences
	list1, isList1 := data1.([]interface{})
	list2, isList2 := data2.([]interface{})
	if isList1 && isList2 {
		for i := 0; i < max(len(list1), len(list2)); i++ {
			var item1, item2 interface{}
			if i < len(list1) {
				item1 = list1[i]
			}
			if i < len(list2) {
				item2 = list2[i]
			}

			itemPath := fmt.Sprintf("%s[%d]", parentPath, i)
			if !reflect.DeepEqual(item1, item2) {
				if isLeafNode(item1) && isLeafNode(item2) {
					fmt.Printf("Found difference on [%s]\n", itemPath)
					printDifferenceContext(red(fmt.Sprintf("- %v", item1)), green(fmt.Sprintf("+ %v", item2)), indent)
				} else {
					compareYAMLWithContext(item1, item2, indent+"  ", itemPath)
				}
			}
		}
	}
}

// Display the context lines around the differences
func printDifferenceContext(diff1, diff2, indent string) {
	fmt.Printf("%s\n", diff1)
	fmt.Printf("%s\n", diff2)
}

// Format YAML values with appropriate indentation
func formatYAMLValue(value interface{}, indent string) string {
	var buffer bytes.Buffer
	encoder := yaml.NewEncoder(&buffer)
	encoder.SetIndent(2)
	_ = encoder.Encode(value)

	lines := strings.Split(buffer.String(), "\n")
	for i := range lines {
		lines[i] = indent + lines[i]
	}
	return strings.Join(lines, "\n")
}

// Determine if a node is a leaf (scalar value)
func isLeafNode(value interface{}) bool {
	_, isMap := value.(map[string]interface{})
	_, isSlice := value.([]interface{})
	return !isMap && !isSlice
}

// Get the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
