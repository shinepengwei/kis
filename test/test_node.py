# -*- coding: utf-8 -*-
import requests
import json
print requests.get('http://localhost:8080/v1/pods').json()

r = requests.post('http://localhost:8080/v1/pods',json ={
	'Name':'podstest4',
	})
print r.text

print requests.get('http://localhost:8080/v1/pods').json()
