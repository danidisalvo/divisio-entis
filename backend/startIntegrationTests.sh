#!/bin/bash

killall divisio-entis-backend

tearDown() {
  killall divisio-entis-backend
  echo "$2"
  exit $1
}

echo "Starting the divisio-entis-backend"
go build -o divisio-entis-backend main.go
./divisio-entis-backend &
sleep 10

#######################################################################################################################

echo
echo "Checking the service's health status"
response=$(curl -s -w "%{http_code}" http://localhost:8080/apis/health --output /dev/null)
if [ $response != 200 ]; then
  tearDown 1 "Health check failed"
fi
echo "Health check ok"

#######################################################################################################################

echo
echo "Adding a node to the root"
response=$(curl -s -w "%{http_code}" -H 'Content-Type: application/json' \
  --data '{"id":"1","name":"A","color":"#ff0000","type":"LEXEME","children":null}' \
  -X PUT http://localhost:8080/apis/nodes --output output.json)
if [ $response != 200 ]; then
  tearDown 1 "Failed to add node A to the root"
fi

VAR=$(curl -s http://localhost:8080/apis/graph | jq  -r '.children[0].name')
if [ "$VAR" != "A" ]; then
  echo "$VAR"
  tearDown 1 "Failed to add node A to the root"
fi
echo "Added node A to the root"

#######################################################################################################################

#######################################################################################################################

echo
echo "Adding duplicated node to the root"
response=$(curl -s -w "%{http_code}" -H 'Content-Type: application/json' \
  --data '{"id":"1","name":"A","color":"#ff0000","type":"LEXEME","children":null}' \
  -X PUT http://localhost:8080/apis/nodes --output output.json)
if [ $response != 400 ]; then
  tearDown 1 "The duplicated node A was added to the root"
fi

#######################################################################################################################

echo
echo
echo "Updating node A"
response=$(curl -s -w "%{http_code}" -H 'Content-Type: application/json' \
  --data '{"id":"1","name":"A","color":"#ffffff","type":"LEXEME","children":null}' \
  -X PUT http://localhost:8080/apis/nodes/0 --output output.json)
if [ $response != 200 ]; then
  tearDown 1 "Failed to update node A"
fi

VAR=$(curl -s http://localhost:8080/apis/graph | jq  -r '.children[0].color')
if [ "$VAR" != "#ffffff" ]; then
  tearDown 1 "Failed to update node A"
fi
echo "Node A updated"

#######################################################################################################################

echo
echo
echo "Adding node B to node A"
response=$(curl -s -w "%{http_code}" -H 'Content-Type: application/json' \
  --data '{"id":"1","name":"A","color":"#ffffff","type":"LEXEME","children":[{"id":"2","name":"B","color":"#ff0000","type":"LEXEME","children":null}]}' \
  -X PUT http://localhost:8080/apis/nodes/0 --output output.json)
if [ $response != 200 ]; then
  tearDown 1 "Failed to add node B to node A"
fi

VAR=$(curl -s http://localhost:8080/apis/graph | jq  -r '.children[0].children[0].name')
if [ "$VAR" != "B" ]; then
  tearDown 1 "Failed to add node B to node A"
fi
echo "Added node B to node A"

#######################################################################################################################

echo
echo
echo "Adding a duplicated node to node A"
response=$(curl -s -w "%{http_code}" -H 'Content-Type: application/json' \
  --data '{"id":"1","name":"A","color":"#ff0000","type":"LEXEME","children":[{"id":"0","name":"ens","color":"#ff0000","type":"LEXEME","children":null}]}' \
  -X PUT http://localhost:8080/apis/nodes/0 --output output.json)
if [ $response != 400 ]; then
  tearDown 1 "The duplicated node ens was added to node A"
fi

#######################################################################################################################

echo
echo
echo "Deleting node B from node A"
response=$(curl -s -w "%{http_code}" -X DELETE http://localhost:8080/apis/nodes/1/2 --output output.json)
if [ $response != 200 ]; then
  tearDown 1 "Failed to delete node B from node A"
fi

VAR=$(curl -s http://localhost:8080/apis/graph | jq  -r '.children[0].children[0].name')
if [ "$VAR" = "B" ]; then
  tearDown 1  "Failed to delete node B from node A"
fi
echo "Deleted node B from node A"

tearDown 0 "test cases succeeded"
