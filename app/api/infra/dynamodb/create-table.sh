#!/bin/bash

set -e

sleep 3
aws dynamodb create-table --cli-input-json file://table.json --endpoint-url http://dynamodb-local:8000
