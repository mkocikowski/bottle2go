package main

import (
	"log"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
// 	"text/template"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	http.HandleFunc("/", root)
}

func error400(w http.ResponseWriter, err error) {
	log.Println(err)
	http.Error(w, "400 Bad Request: " + err.Error(), 400)
}

func parse(body []byte) (map[string]interface{}, error) {

	data := map[string]interface{}{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	for _, field := range []string{"name", "age"} {
		if _, ok := data[field]; !ok {
			return nil, fmt.Errorf("'%s' is required", field)
		}
	}

	switch data["age"].(type) {
	case float64:
		break
	default:
		return nil, fmt.Errorf("'age' must be a number")
	}

	return data, nil
}

func root(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain")
	body, _ := ioutil.ReadAll(r.Body)

	data, err := parse(body)
	if err != nil {
		error400(w, err)
		return
	}

	s := fmt.Sprintf("%s's age is %d\n", data["name"], int(data["age"].(float64)))
	fmt.Fprint(w, s)
}

func main() {
	err := http.ListenAndServe("127.0.0.1:8082", nil)
	if err != nil {
		log.Fatal(err)
	}
}
