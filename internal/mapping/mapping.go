package mapping

import (
	"encoding/json"
	"log"
)

func CreateRequest(mappingMap map[string]interface{}) ([]byte, error) {

	var request map[string]interface{}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}

	return jsonRequest, nil
}
