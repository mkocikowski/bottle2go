package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func handle(w http.ResponseWriter, r *http.Request) error {

	// this is like python 'finally'
	defer r.Body.Close()

	// anonymous struct: https://talks.golang.org/2012/10things.slide#2
	var data struct {
		Name string
		Age  int64
	}

	// using struct for data ensures that if Name is not a string or age
	// is not a number there is an error
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return err
	}

	// 'zero' value (value set by go on creation) for string is "", for
	// float64 is 0, so when data struct is created, its Name == "" and
	// Age == 0.
	if data.Name == "" || data.Age == 0 {
		return fmt.Errorf("name and age are required")
	}

	fmt.Fprintf(w, "%v's age is %v\n", data.Name, data.Age)
	return nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	http.HandleFunc("/", errorHandler(handle))
	err := http.ListenAndServe("127.0.0.1:8082", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// errorHandler handles the errors returned by an http handler,
// almost like a Python decorator.
func errorHandler(h func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("bad request: %v\n", err)
		}
	}
}
