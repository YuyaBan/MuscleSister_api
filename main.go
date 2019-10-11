package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

type TODO struct {
	ID   int    "json:id"
	Name string "json:name"
	User string "json:user"
	Done bool   "json:done"
}

var (
	mu    = sync.RWMutex{}
	cache = make(map[int]TODO)
	maxID = int32(0)
)

func Road(w http.ResponseWriter, r *http.Request) {
	log.Println("[+] Get start")

	idStr := strings.TrimPrefix(r.RequestURI, "/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("[-] Get Faild1")
		return
	}

	mu.RLock()
	data, ok := cache[id]
	mu.RUnlock()
	if !ok {
		log.Println("[-] Get Failed2")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println("[-] Get Failed3")
		log.Println(err)
	}

	log.Println("[+] Get Success")
}

func Create(w http.ResponseWriter, r *http.Request) {
	var todo TODO
	log.Println("[+] Create start")

	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		log.Println("[-] Create Faild1")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if todo.ID == 0 {
		mu.Lock()
		todo.ID = int(maxID)
		maxID++
		mu.Unlock()
	}

	mu.Lock()
	cache[todo.ID] = todo
	mu.Unlock()

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(todo)
	log.Println("[+] Post Success")
}

func Update(w http.ResponseWriter, r *http.Request) {
	log.Println("[+] Update Start")

	idStr := strings.TrimPrefix(r.RequestURI, "/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("[-] Update Failed1")
		return
	}

	mu.Lock()
	todo, ok := cache[id]
	mu.Unlock()
	if !ok {
		log.Println("[-] Update Failed2")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	todo.Done = true

	mu.Lock()
	cache[id] = todo
	mu.Unlock()

	w.WriteHeader(http.StatusAccepted)

	json.NewEncoder(w).Encode(todo)
	log.Println("[+] Update Success")
}

func Delete(w http.ResponseWriter, r *http.Request) {
	log.Println("[+] Delete start")

	idStr := strings.TrimPrefix(r.RequestURI, "/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("[-] Delete Faild1")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mu.Lock()
	delete(cache, id)
	mu.Unlock()

	w.WriteHeader(http.StatusAccepted)
	log.Println("[+] Delete Success")

}

func main() {
	port, _ := strconv.Atoi(os.Args[1])
	// envPort := os.Getenv("PORT")
	// if len(envPort) != 0 {
	// 	port = envPort
	// }

	mux := http.DefaultServeMux
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			Road(w, r)
		case http.MethodPost:
			Create(w, r)
		case http.MethodPatch, http.MethodPut:
			Update(w, r)
		case http.MethodDelete:
			Delete(w, r)
		}
	})

	http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	// http.ListenAndServe(port, mux)
}
