#!/bin/bash

killall divisio-entis-backend
rm volume/graph.json

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
echo "Resetting the graph"
response=$(curl -s -w "%{http_code}" -X DELETE http://localhost:8080/apis/graph)
if [ $response != 204 ]; then
  tearDown 1 "Failed to delete the graph"
fi

#######################################################################################################################

echo
echo "Adding a node to the root"
response=$(curl -s -w "%{http_code}" -H 'Content-Type: application/json' \
  --data '{"id":"1","name":"A","color":"#ff0000","type":"lexeme","children":null}' \
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

echo
echo "Adding duplicated node to the root"
response=$(curl -s -w "%{http_code}" -H 'Content-Type: application/json' \
  --data '{"id":"1","name":"A","color":"#ff0000","type":"lexeme","children":null}' \
  -X PUT http://localhost:8080/apis/nodes --output output.json)
if [ $response != 400 ]; then
  tearDown 1 "The duplicated node A was added to the root"
fi

#######################################################################################################################

echo
echo
echo "Updating node A"
response=$(curl -s -w "%{http_code}" -H 'Content-Type: application/json' \
  --data '{"id":"1","name":"A","color":"#ffffff","type":"lexeme","children":null}' \
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
  --data '{"id":"1","name":"AAA","color":"#ffffff","type":"division","children":[{"id":"2","name":"B","color":"#ff0000","type":"lexeme","children":null}]}' \
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
echo "Renaming node A"
response=$(curl -s -w "%{http_code}" -H 'Content-Type: application/json' \
  --data '{"id":"1","name":"AAA","color":"#ffffff","type":"division","children":null}' \
  -X PUT http://localhost:8080/apis/nodes/0 --output output.json)
if [ $response != 200 ]; then
  tearDown 1 "Failed to update node A"
fi

VAR=$(curl -s http://localhost:8080/apis/graph | jq  -r '.children[0].name')
if [ "$VAR" != "AAA" ]; then
  tearDown 1 "Failed to update node A's name"
fi

VAR=$(curl -s http://localhost:8080/apis/graph | jq  -r '.children[0].type')
if [ "$VAR" != "division" ]; then
  tearDown 1 "Failed to update node A's division"
fi

VAR=$(curl -s http://localhost:8080/apis/graph | jq  -r '.children[0].children[0].name')
if [ "$VAR" != "B" ]; then
  tearDown 1 "The child of node A has been deleted"
fi

echo "Node A updated to node AAA"

#######################################################################################################################

echo
echo
echo "Adding a duplicated node to node AAA"
response=$(curl -s -w "%{http_code}" -H 'Content-Type: application/json' \
  --data '{"id":"1","name":"AAA","color":"#ffffff","type":"division","children":[{"id":"0","name":"ens","color":"#ff0000","type":"lexeme","children":null}]}' \
  -X PUT http://localhost:8080/apis/nodes/0 --output output.json)
if [ $response != 400 ]; then
  tearDown 1 "The duplicated node ens was added to node A"
fi

#######################################################################################################################

echo
echo
echo "Deleting node B from node AAAA"
response=$(curl -s -w "%{http_code}" -X DELETE http://localhost:8080/apis/nodes/1/2 --output output.json)
if [ $response != 200 ]; then
  tearDown 1 "Failed to delete node B from node A"
fi

VAR=$(curl -s http://localhost:8080/apis/graph | jq  -r '.children[0].children[0].name')
if [ "$VAR" = "B" ]; then
  tearDown 1  "Failed to delete node B from node A"
fi
echo "Deleted node B from node A"

#######################################################################################################################

echo
echo
echo "Adding node C to node AAA"
response=$(curl -s -w "%{http_code}" -H 'Content-Type: application/json' \
  --data '{"id":"1","name":"AAA","color":"#ffffff","type":"division","children":[{"id":"3","name":"C","color":"#ff0000","type":"lexeme","children":null}]}' \
  -X PUT http://localhost:8080/apis/nodes/0 --output output.json)
if [ $response != 200 ]; then
  tearDown 1 "Failed to add node C to node AAA"
fi

VAR=$(curl -s http://localhost:8080/apis/graph | jq  -r '.children[0].children[0].name')
if [ "$VAR" != "C" ]; then
  tearDown 1 "Failed to add node C to node A"
fi
echo "Added node C to node A"

#######################################################################################################################

#######################################################################################################################

echo
echo
echo "Adding node D to node C"
response=$(curl -s -w "%{http_code}" -H 'Content-Type: application/json' \
  --data '{"id":"3","name":"C","color":"#ff0000","type":"lexeme","children":[{"id":"4","name":"D","color":"#ff0000","type":"lexeme","children":null}]}' \
  -X PUT http://localhost:8080/apis/nodes/1 --output output.json)
if [ $response != 200 ]; then
  tearDown 1 "Failed to add node D to node C"
fi

VAR=$(curl -s http://localhost:8080/apis/graph | jq  -r '.children[0].children[0].name')
if [ "$VAR" != "C" ]; then
  tearDown 1 "Failed to add node D to node C"
fi
echo "Added node D to node C"

#######################################################################################################################

echo
echo
echo "Getting C's target nodes"
response=$(curl -s -w "%{http_code}" -H 'Content-Type: application/json' \
  -X GET http://localhost:8080/apis/nodes/3/targets --output output.json)

cat output.json

if [ $response != 200 ]; then
  tearDown 1 "Failed to get C's target noded"
fi

#######################################################################################################################

echo
echo
echo "Moving node C to the root"
response=$(curl -s -w "%{http_code}" -H 'Content-Type: application/json' \
  -X POST http://localhost:8080/apis/nodes/1/3/0 --output output.json)
if [ $response != 200 ]; then
  tearDown 1 "Failed to move node H to the root"
fi

VAR=$(curl -s http://localhost:8080/apis/graph | jq  -r '.children[1].name')
if [ "$VAR" != "C" ]; then
  tearDown 1 "Failed to move node C to the root"
fi
echo "Moved node C to the root"

tearDown 0 "test cases succeeded"

