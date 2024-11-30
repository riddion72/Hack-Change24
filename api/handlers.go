package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	lr "main/internal/schem_reader"
)

var mlList = []string{
	"examples/model.json",
}

var dbList = []string{
	"examples/credit_form_schem.json",
	"examples/education_departmen_schem.json",
	"examples/fedresource_schem.json",
	"examples/UCB_schem.json",
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func getDBlist(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	for _, d := range dbList {
		mlS, err := lr.ReadJSONFromFile(d)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error reading file: %v\n", err)
			log.Printf("Error reading file: %v\n", err)
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(mlS)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func getMLlist(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	for _, d := range mlList {
		mlS, err := lr.ReadJSONFromFile(d)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error reading file: %v\n", err)
			log.Printf("Error reading file: %v\n", err)
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(mlS)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func createPattern(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("ошибка при разборе JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	path := fmt.Sprintf("examples/%s.json", result["name"])

	file, err := os.Create(path)
	if err != nil {
		log.Printf("Ошибка при создании файла: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()
	_, err = file.Write(body)
	if err != nil {
		log.Printf("Ошибка при записи данных: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Printf("Received: %v\n", result)
}

func performRequest(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("ошибка при разборе JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	path := fmt.Sprintf("examples/%s.json", result["name"])
	patternFile, err := os.Open(path)
	if err != nil {
		log.Printf("Ошибка при открытии файла: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer patternFile.Close()
	patternBody, err := lr.ReadJSONFromFile(patternFile.Name())
	var pattern map[string]interface{}
	err = json.Unmarshal(patternBody, &pattern)
	if err != nil {
		log.Printf("ошибка при разборе JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Print(pattern)
	w.WriteHeader(http.StatusOK)
}
