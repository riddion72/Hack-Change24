package mapping

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	conv "main/internal/convertors"
)

// var dbList = []string{
// 	"credit_form_schem",
// 	"education_departmen_schem",
// 	"fedresource_schem",
// 	"UCB_schem",
// }

func CreateRequest(mappingMap map[string]interface{}) ([]byte, error) {

	forReplace := make(map[string]string)

	prepareValues(mappingMap, forReplace)

	jsonRequest, err := json.Marshal(mappingMap)
	if err != nil {
		log.Print(err)

		return nil, err
	}
	for k, v := range forReplace {
		jsonRequest = bytes.ReplaceAll(jsonRequest, []byte(k), []byte(v))
	}

	return jsonRequest, nil
}

func prepareValues(value interface{}, answer map[string]string) {
	switch v := value.(type) {
	case map[string]interface{}:
		for _, val := range v {
			if _, ok := val.(map[string]interface{}); !ok {
				key, data := mapping(fmt.Sprint(val))
				answer[key] = data
			} else {
				prepareValues(val, answer)
			}
		}
	case []interface{}:
		for _, val := range v {
			prepareValues(val, answer)
		}
	}
}

func mapping(instruction string) (string, string) {
	parts := strings.Split(instruction, " ")
	if len(parts) < 2 {
		return instruction, instruction
	}
	path := fmt.Sprintf("examples/%s.json", parts[0])
	dbFile, err := os.Open(path)
	if err != nil {
		log.Printf("Ошибка при открытии файла: %v", err)
		return "", ""
	}
	defer dbFile.Close()

	byteValue, err := io.ReadAll(dbFile)
	if err != nil {
		return "", ""
	}
	var promt string
	for i := 3; i < len(parts); i++ {
		promt += " " + parts[i]
	}

	return instruction, applyFunction(parts[2], extractValue(byteValue, parts[1]), promt)
}

func extractValue(dbAnswer []byte, requiredFieldName string) string {
	var data map[string]interface{}
	err := json.Unmarshal(dbAnswer, &data)
	if err != nil {
		log.Print(err)
		return ""
	}
	return findValues(data, requiredFieldName)
}

func findValues(value interface{}, requiredFieldName string) string {
	var requiredFieldValue string
	switch v := value.(type) {
	case map[string]interface{}:
		for key, val := range v {
			if key == requiredFieldName {
				requiredFieldValue += fmt.Sprintf("%v", val)
			}
			if _, ok := val.(map[string]interface{}); ok {
				requiredFieldValue += findValues(val, requiredFieldName)
			}
		}
	case []interface{}:
		for _, val := range v {
			requiredFieldValue += findValues(val, requiredFieldName)
		}
	}
	return requiredFieldValue
}

func applyFunction(funcName string, funcArgs string, promt string) string {
	switch funcName {
	case "calcAge":
		return conv.CalculateAge(funcArgs)
	case "average":
		return conv.Average(funcArgs)
	case "calcItems":
		return conv.CalcItems(funcArgs)
	case "askAI":
		return conv.AskAI(funcArgs + " " + promt)
	case "insert":
		return conv.Insert(funcArgs)
	default:
		return funcArgs
	}
}
