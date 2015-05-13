RESTful app in Python comapred to one written in Go
---

The purpose of this is to show how a RESTful app written in Python
(using Bottle) translates into Go. The app accepts `POST /`, with the
body being a json object `{"name":"Monkey","age":10}`. Input is
validated; on validation failure status 400 is returned with body
describing the error, on success 200 and the body is a text/plain
sentence `Monkey's age is 10\n`. Sample input and output:


	$ curl -i -XPOST 127.0.0.1:8081/ -d '{"name":"Monkey","age":10}'

	HTTP/1.0 200 OK
	Date: Wed, 13 May 2015 20:07:55 GMT
	Server: WSGIServer/0.1 Python/2.7.6
	Content-Length: 19
	Content-Type: text/plain

	Monkey's age is 10

	$ curl -i -XPOST 127.0.0.1:8081/ -d '{"name":"Monkey"}'

	HTTP/1.0 400 Bad Request
	Date: Wed, 13 May 2015 20:59:13 GMT
	Server: WSGIServer/0.1 Python/2.7.6
	Content-Length: 35
	Content-Type: text/html; charset=UTF-8

	400 Bad Request: 'age' is required


Dependencies
------------
Python needs (`pip install bottle WebTest`):
```
astroid==1.3.6
beautifulsoup4==4.3.2
bottle==0.12.8
logilab-common==0.63.2
six==1.9.0
waitress==0.8.9
WebOb==1.4.1
WebTest==2.0.18
```
Go uses only standard library.


Logging
-------
```python
logger = logging.getLogger(__name__)
logging.basicConfig(
  level=logging.DEBUG,
  format='%(levelname)s %(filename)s:%(funcName)s:%(lineno)d %(message)s'
)
```
```go
log.SetFlags(log.LstdFlags | log.Llongfile)
```


Routes
------
```python
@app.post('/')
def handle():
  ...
```
```go
http.HandleFunc("/", root)
```


Errors
------
```python
@app.error(400)
def error400(err):
    logger.error(err)
    return "400 Bad Request: %s\n" % err.body
...
raise bottle.HTTPError(status=400, body=exc)
```
```go
func error400(w http.ResponseWriter, err error) {
	log.Println(err)
	http.Error(w, "400 Bad Request: "+err.Error(), 400)
}
...
if err != nil {
	error400(w, err)
	return
}
```


Handle request
--------------
```python
@app.post('/')
def handle():
    try:
        bottle.response.set_header("Content-Type", "text/plain")
        body = bottle.request.body.read()
        data = parse(body)
        s = "%s's age is %s\n" % (data["name"], int(data["age"]))
        return s
    except (StandardError,) as exc:
        raise bottle.HTTPError(status=400, body=exc)
```
```go
func handle(w http.ResponseWriter, r *http.Request) {
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
```

Parse request
-------------
```python
def parse(body):
    data = json.loads(body)
    for field in ["name", "age"]:
        if field not in data:
            raise ValueError("'%s' is required" % field)
    if not isinstance(data["age"], (int, float)):
        raise ValueError("'age' must be a number")
    return data
```
```go
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
```


Run Tests
---------
```shell
cd bottle2go
python srv_test.py
go test ./...
```


```python
app = webtest.TestApp(srv.app)
...
tests = [
    {"in": '{"name":"Monkey","age":10}', "out": "Monkey's age is 10\n", "code": 200},
    ...
]
for test in tests:
    response = app.post('/', test["in"], status=test["code"])
    self.assertEqual(response.body, test["out"])
    self.assertEqual(response.status_int, test["code"])
```
```go
ts := httptest.NewServer(nil)
defer ts.Close()
...
tests := []struct{
	in string
	out string
	status int
}{
	{`{"name":"Monkey","age":10}`, "Monkey's age is 10\n", 200},
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
```
