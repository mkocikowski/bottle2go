RESTful app in Python translated into Go
---

This branch has the idiomatic go server. The 'main' branch has as close
a translation of python to go as possible.

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

Run
---
Python runs on `127.0.0.1:8081`:
```shell
python srv.py
```
Go runs on `127.0.0.1:8082`:
```shell
go run srv.go
```

Run Tests
---------
```shell
cd bottle2go
python srv_test.py
go test ./...
```

