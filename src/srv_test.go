package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
)

func TestRoot(t *testing.T) {

	ts := httptest.NewServer(nil)
	defer ts.Close()

	tests := []struct{
		in string
		out string
		status int
	}{
		{`{"name":"Monkey","age":10}`, "Monkey's age is 10\n", 200},
		{`{"name":"Monkey","age":10.1}`, "Monkey's age is 10\n", 200},
		{`{"name":"Monkey","age":"10"}`, "400 Bad Request: 'age' must be a number\n", 400},
		{`{"name":"Monkey"}`, "400 Bad Request: 'age' is required\n", 400},
		{"", "400 Bad Request: unexpected end of JSON input\n", 400},
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
