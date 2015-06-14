package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func handle(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()

	var data struct {
		Name string
		Age  float64
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return err
	}
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
			log.Println("bad request: %v", err)
		}
	}
}
