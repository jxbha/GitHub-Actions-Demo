package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

var DOC string = "tmp.json"

type Document struct {
	LastId     int
	Collection []Widget
}
type Widget struct {
	Id    int
	Key   string `json:"key"`
	Value string `json:"value"`
}

func createWidget(w http.ResponseWriter, req *http.Request) {
	var widget Widget
	var document Document
	body, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = json.Unmarshal(body, &widget)
	if err != nil {
		fmt.Println(err.Error())
	}
	f, err := os.ReadFile(DOC)
	if err != nil {
		fmt.Println("File does not exist")
		os.Create(DOC)
	}
	err = json.Unmarshal(f, &document)
	if err != nil {
		fmt.Println(err.Error())
	}
	widget.Id = document.LastId + 1
	document.LastId += 1

	document.Collection = append(document.Collection, widget)
	result, err := json.Marshal(document)

	if err != nil {
		fmt.Println(err.Error())
	}
	os.WriteFile(DOC, result, 0600)
	fmt.Fprint(w, "success; created")

}
func getWidget(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.PathValue("id"))
	if err != nil {
		fmt.Println("problem reading request")
	}
	f, err := os.ReadFile(DOC)
	if err != nil {
		fmt.Println("problem reading file")
	}
	var document Document
	json.Unmarshal(f, &document)

	m := make(map[int]Widget)
	for _, v := range document.Collection {
		m[v.Id] = v
	}
	if _, ok := m[id]; !ok {
		fmt.Fprint(w, "record not found")
		return
	}
	result, err := json.MarshalIndent(m[id], "", "  ")
	fmt.Fprint(w, string(result))
}
func updateWidget(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.PathValue("id"))
	if err != nil {
		fmt.Println("problem reading request")
	}
	body, err := io.ReadAll(req.Body)

	f, err := os.ReadFile(DOC)
	if err != nil {
		fmt.Println("problem reading file")
	}
	var document Document
	json.Unmarshal(f, &document)

	var widgets []Widget
	m := make(map[int]Widget)
	for _, v := range document.Collection {
		m[v.Id] = v
	}
	if _, ok := m[id]; !ok {
		fmt.Fprint(w, "record not found")
		return
	}
	var widget Widget
	json.Unmarshal(body, &widget)
	widget.Id = id
	m[id] = widget
	for _, v := range m {
		widgets = append(widgets, v)
	}
	document.Collection = widgets

	result, err := json.Marshal(document)
	if err != nil {
		fmt.Println(err.Error())
	}
	os.WriteFile(DOC, result, 0600)
	fmt.Fprint(w, "Success; updated")
}
func deleteWidget(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.PathValue("id"))
	if err != nil {
		fmt.Println("problem reading request")
	}

	f, err := os.ReadFile(DOC)
	if err != nil {
		fmt.Println("problem reading file")
	}

	var document Document
	json.Unmarshal(f, &document)

	m := make(map[int]Widget)
	for _, v := range document.Collection {
		m[v.Id] = v
	}

	if _, ok := m[id]; !ok {
		fmt.Fprint(w, "record not found")
		return
	}
	delete(m, id)

	var widgets []Widget
	for _, v := range m {
		widgets = append(widgets, v)
	}

	result, err := json.Marshal(document)
	if err != nil {
		fmt.Println(err.Error())
	}

	os.WriteFile(DOC, result, 0660)
	fmt.Fprint(w, "success; deleted")
}

func listWidget(w http.ResponseWriter, req *http.Request) {
	var document Document
	f, _ := os.ReadFile(DOC)
	json.Unmarshal(f, &document)
	result, _ := json.MarshalIndent(document, "", "  ")
	fmt.Fprint(w, string(result))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", listWidget)
	mux.HandleFunc("/create", createWidget)
	mux.HandleFunc("/update/{id}", updateWidget)
	mux.HandleFunc("/delete/{id}", deleteWidget)
	mux.HandleFunc("/{id}", getWidget)

	port := ":4040"
	fmt.Printf("listening on %s...\n", port)
	http.ListenAndServe(port, mux)
}
