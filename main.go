package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type Task struct {
	Name      string `json:"name"`
	Priority  string `json:"priority"`
	DateAdded string `json:"date_added"`
}

var tasks = []Task{}
var mu sync.Mutex

func fetchTasksFromPythonAPI() error {
	resp, err := http.Get("http://localhost:5000/tasks")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result map[string][]Task
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	mu.Lock()
	tasks = result["all"]
	mu.Unlock()
	return nil
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	if err := fetchTasksFromPythonAPI(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(tasks)
}

func completeTask(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call Python API to mark task as completed
	jsonTask, _ := json.Marshal(task)
	resp, err := http.Post("http://localhost:5000/tasks/complete", "application/json", bytes.NewBuffer(jsonTask))
	if err != nil || resp.StatusCode != http.StatusNoContent {
		http.Error(w, "Failed to complete task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	// Serve static files from the "static" directory.
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// Define endpoints for task operations.
	http.HandleFunc("/tasks", getTasks)
	http.HandleFunc("/tasks/complete", completeTask)

	// Start the server on port 8080.
	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
