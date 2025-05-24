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
	Comment   string `json:"comment,omitempty"`
}

type PassFail struct {
	Date string `json:"date"`
	Pass int    `json:"pass"`
	Fail int    `json:"fail"`
}

type TurnAroundTime struct {
	Date           string  `json:"date"`
	TurnAroundTime float64 `json:"turn_around_time"`
}

type OpenCloseMonthly struct {
	Date  string `json:"date"`
	Open  int    `json:"open"`
	Close int    `json:"close"`
}

type TypeData struct {
	Type string `json:"type"`
	Qty  int    `json:"qty"`
}

type TasksCompletedThisWeek struct {
	Count int `json:"count"`
}

var tasks = []Task{}
var mu sync.Mutex

func fetchTasksFromPythonAPI() error {
	resp, err := http.Get("http://localhost:5000/tasks")
	if err != nil {
		log.Printf("Error fetching tasks from Python API: %v", err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading tasks response: %v", err)
		return err
	}

	log.Printf("Raw tasks response from Python API: %s", string(body))

	var result map[string][]Task
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Error unmarshaling tasks JSON: %v", err)
		return err
	}

	mu.Lock()
	tasks = result["all"]
	log.Printf("Parsed tasks: %v", tasks)
	mu.Unlock()
	return nil
}

func fetchPassFailFromPythonAPI() (PassFail, error) {
	resp, err := http.Get("http://localhost:5000/pass_fail")
	if err != nil {
		return PassFail{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return PassFail{}, err
	}

	var passFail PassFail
	if err := json.Unmarshal(body, &passFail); err != nil {
		return PassFail{}, err
	}
	return passFail, nil
}

func fetchTurnAroundTimeFromPythonAPI() ([]TurnAroundTime, error) {
	resp, err := http.Get("http://localhost:5000/turn_around_time")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var turnAroundTimes []TurnAroundTime
	if err := json.Unmarshal(body, &turnAroundTimes); err != nil {
		return nil, err
	}
	return turnAroundTimes, nil
}

func fetchOpenCloseMonthlyFromPythonAPI() ([]OpenCloseMonthly, error) {
	resp, err := http.Get("http://localhost:5000/open_close_monthly")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var openCloseData []OpenCloseMonthly
	if err := json.Unmarshal(body, &openCloseData); err != nil {
		return nil, err
	}
	return openCloseData, nil
}

func fetchTypesFromPythonAPI() ([]TypeData, error) {
	resp, err := http.Get("http://localhost:5000/types")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var typesData []TypeData
	if err := json.Unmarshal(body, &typesData); err != nil {
		return nil, err
	}
	return typesData, nil
}

func fetchTasksCompletedThisWeekFromPythonAPI() (TasksCompletedThisWeek, error) {
	resp, err := http.Get("http://localhost:5000/tasks_completed_this_week")
	if err != nil {
		return TasksCompletedThisWeek{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return TasksCompletedThisWeek{}, err
	}

	var tasksCompleted TasksCompletedThisWeek
	if err := json.Unmarshal(body, &tasksCompleted); err != nil {
		return TasksCompletedThisWeek{}, err
	}
	return tasksCompleted, nil
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	if err := fetchTasksFromPythonAPI(); err != nil {
		log.Printf("Error in getTasks: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	log.Printf("Sending tasks to client: %v", tasks)
	json.NewEncoder(w).Encode(tasks)
}

func completeTask(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var task map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonTask, _ := json.Marshal(task)
	resp, err := http.Post(
		"http://localhost:5000/tasks/complete",
		"application/json",
		bytes.NewBuffer(jsonTask),
	)
	if err != nil || resp.StatusCode != http.StatusNoContent {
		http.Error(w, "Failed to complete task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getPassFail(w http.ResponseWriter, r *http.Request) {
	passFail, err := fetchPassFailFromPythonAPI()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(passFail)
}

func getTurnAroundTime(w http.ResponseWriter, r *http.Request) {
	turnAroundTimes, err := fetchTurnAroundTimeFromPythonAPI()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(turnAroundTimes)
}

func getOpenCloseMonthly(w http.ResponseWriter, r *http.Request) {
	openCloseData, err := fetchOpenCloseMonthlyFromPythonAPI()
	if err != nil {
		log.Printf("Error in getOpenCloseMonthly: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(openCloseData)
}

func getTypes(w http.ResponseWriter, r *http.Request) {
	typesData, err := fetchTypesFromPythonAPI()
	if err != nil {
		log.Printf("Error in getTypes: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(typesData)
}

func getTasksCompletedThisWeek(w http.ResponseWriter, r *http.Request) {
	tasksCompleted, err := fetchTasksCompletedThisWeekFromPythonAPI()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(tasksCompleted)
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/tasks", getTasks)
	http.HandleFunc("/tasks/complete", completeTask)
	http.HandleFunc("/pass_fail", getPassFail)
	http.HandleFunc("/turn_around_time", getTurnAroundTime)
	http.HandleFunc("/open_close_monthly", getOpenCloseMonthly)
	http.HandleFunc("/types", getTypes)
	http.HandleFunc("/tasks_completed_this_week", getTasksCompletedThisWeek)
	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
