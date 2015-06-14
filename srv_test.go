package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
)

func TestHandle(t *testing.T) {

	ts := httptest.NewServer(errorHandler(handle))
	defer ts.Close()

	tests := []struct{
		in string
		out string
		status int
	}{
		{`{"name":"Monkey","age":10}`, "Monkey's age is 10\n", 200},
		{`{"name":"Monkey","age":10.1}`, "json: cannot unmarshal number 10.1 into Go value of type int64\n", 400},
		{`{"name":"Monkey","age":"10"}`, "json: cannot unmarshal string into Go value of type int64\n", 400},
		{`{"name":"Monkey"}`, "name and age are required\n", 400},
		{"", "EOF\n", 400},
	}

	for _, test := range tests   {

		res, err := http.Post(ts.URL, "application/json", strings.NewReader(test.in))
		if err != nil {
			t.Fatal(err)
		}

		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		if string(body) != test.out {
			t.Fatalf("expected '%s', got '%s'", test.out, string(body))
		}

		if res.StatusCode != test.status {
			t.Fatalf("expected code %d, got %d", test.status, res.StatusCode)
		}

	}

}
