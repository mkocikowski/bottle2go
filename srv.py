# -*- coding: UTF-8 -*-

import logging
import json

import bottle

logger = logging.getLogger(__name__)
app = bottle.Bottle()

@app.error(400)
def error400(err):
    logger.error(err)
    return "400 Bad Request: %s\n" % err.body


def parse(body):

    data = json.loads(body)

    for field in ["name", "age"]:
        if field not in data:
            raise ValueError("'%s' is required" % field)

    if not isinstance(data["age"], (int, float)):
        raise ValueError("'age' must be a number")

    return data


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


if __name__ == '__main__':
    logging.basicConfig(
        level=logging.DEBUG,
        format='%(levelname)s %(filename)s:%(funcName)s:%(lineno)d %(message)s'
    )
    bottle.run(app=app, host="127.0.0.1", port=8081, reloader=True)

# curl -i -XPOST localhost:8081/ -d '{"name":"Monkey","age":10}'
