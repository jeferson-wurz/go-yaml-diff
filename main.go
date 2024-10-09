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
	// Imprime o cabeçalho estilizado do programa
	printHeader()

	// Define os nomes dos arquivos YAML a serem comparados
	file1Name := "examples/file1.yaml"
	file2Name := "examples/file2.yaml"

	// Lê o conteúdo dos arquivos YAML
	file1Content, err := readFile(file1Name)
	if err != nil {
		log.Fatalf("Erro ao ler o arquivo %s: %v", file1Name, err)
	}

	file2Content, err := readFile(file2Name)
	if err != nil {
		log.Fatalf("Erro ao ler o arquivo %s: %v", file2Name, err)
	}

	// Converte o conteúdo dos arquivos em documentos YAML
	docs1 := parseYAMLDocuments(file1Content)
	docs2 := parseYAMLDocuments(file2Content)

	// Cria um mapeamento dos documentos por "kind" e "metadata.name"
	map1 := mapYAMLDocuments(docs1)
	map2 := mapYAMLDocuments(docs2)

	// Compara os documentos YAML e exibe as diferenças encontradas
	fmt.Printf("\nDifferences found between the YAML files:\n\n")
	compareYAMLMaps(map1, map2)
}

// Função para imprimir o cabeçalho estilizado do programa
func printHeader() {
	blue := color.New(color.FgBlue, color.Bold).SprintFunc()
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	fmt.Println(blue("========================================"))
	fmt.Println(green("             go-yaml-diff               "))
	fmt.Println(blue("========================================"))
}

// Função para ler o conteúdo de um arquivo
func readFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Função para analisar múltiplos documentos YAML em uma lista de interfaces
func parseYAMLDocuments(content string) []interface{} {
	var docs []interface{}
	decoder := yaml.NewDecoder(strings.NewReader(content))
	for {
		var doc interface{}
		if err := decoder.Decode(&doc); err != nil {
			break
		}
		docs = append(docs, doc)
	}
	return docs
}

// Função para mapear documentos YAML por "kind" e "metadata.name"
func mapYAMLDocuments(docs []interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, doc := range docs {
		if obj, ok := doc.(map[string]interface{}); ok {
			kind := obj["kind"]
			metadata, hasMetadata := obj["metadata"].(map[string]interface{})
			name := ""
			if hasMetadata {
				name = metadata["name"].(string)
			}
			if kind != nil && name != "" {
				key := fmt.Sprintf("%s/%s", kind, name)
				result[key] = obj
			}
		}
	}
	return result
}

// Função para comparar dois mapas de documentos YAML e imprimir diferenças com contexto
func compareYAMLMaps(map1, map2 map[string]interface{}) {
	for key, obj1 := range map1 {
		if obj2, exists := map2[key]; exists {
			fmt.Printf("Comparing object %s:\n", key)
			compareYAMLWithContext(obj1, obj2, "", key)
		} else {
			fmt.Printf("Object missing in file 2: %s\n", key)
		}
		fmt.Println("---")
	}

	// Verifica objetos que estão presentes apenas em map2
	for key := range map2 {
		if _, exists := map1[key]; !exists {
			fmt.Printf("Object missing in file 1: %s\n", key)
			fmt.Println("---")
		}
	}
}

// Função para comparar dados YAML e imprimir as diferenças com linhas de contexto, removendo redundâncias
func compareYAMLWithContext(data1, data2 interface{}, indent, parentPath string) {
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	map1, ok1 := data1.(map[string]interface{})
	map2, ok2 := data2.(map[string]interface{})

	// Comparação de mapas (objetos)
	if ok1 && ok2 {
		for key := range map1 {
			value1 := map1[key]
			value2, exists := map2[key]

			currentPath := fmt.Sprintf("%s.%s", parentPath, key)
			if exists && !reflect.DeepEqual(value1, value2) {
				if isLeafNode(value1) && isLeafNode(value2) {
					fmt.Printf("Found difference on [%s]\n", currentPath)
					printContext(red(fmt.Sprintf("- %s: %v", key, value1)), green(fmt.Sprintf("+ %s: %v", key, value2)), indent)
				} else {
					compareYAMLWithContext(value1, value2, indent+"  ", currentPath)
				}
			} else if !exists {
				fmt.Printf("%s%s\n", red(fmt.Sprintf("- %s: %v", key, formatYAMLValue(value1, indent))))
			}
		}
	}

	// Comparação de listas (arrays)
	list1, okList1 := data1.([]interface{})
	list2, okList2 := data2.([]interface{})
	if okList1 && okList2 {
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
					printContext(red(fmt.Sprintf("- %v", item1)), green(fmt.Sprintf("+ %v", item2)), indent)
				} else {
					compareYAMLWithContext(item1, item2, indent+"  ", itemPath)
				}
			}
		}
	}
}

// Função para exibir as linhas de contexto ao redor das diferenças
func printContext(diff1, diff2, indent string) {
	fmt.Printf("%s\n", diff1)
	fmt.Printf("%s\n", diff2)
}

// Função para formatar valores YAML com a indentação apropriada
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

// Função para determinar se um nó é uma folha (valor escalar)
func isLeafNode(value interface{}) bool {
	_, isMap := value.(map[string]interface{})
	_, isSlice := value.([]interface{})
	return !isMap && !isSlice
}

// Função para obter o máximo de dois inteiros
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
