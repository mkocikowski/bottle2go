# -*- coding: UTF-8 -*-

import unittest
import logging
import json

import webtest

import srv


class SrvTest(unittest.TestCase):

    def test_root(self):

        tests = [
            {"in": '{"name":"Monkey","age":10}', "out": "Monkey's age is 10\n", "code": 200},
            {"in": '{"name":"Monkey","age":10.1}', "out": "Monkey's age is 10\n", "code": 200},
            {"in": '{"name":"Monkey","age":"10"}', "out": "400 Bad Request: 'age' must be a number\n", "code": 400},
            {"in": '{"name":"Monkey"}', "out": "400 Bad Request: 'age' is required\n", "code": 400},
            {"in": '', "out": "400 Bad Request: No JSON object could be decoded\n", "code": 400},
        ]

        app = webtest.TestApp(srv.app)

        for test in tests:
            response = app.post('/', test["in"], status=test["code"])
            self.assertEqual(response.body, test["out"])
            self.assertEqual(response.status_int, test["code"])


if __name__ == "__main__":

    logging.basicConfig(level=logging.CRITICAL)
    unittest.main()


