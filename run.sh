#!/bin/bash

echo "Export PATH"
export PATH=$(go env GOPATH)/bin:$PATH

echo "Formatting Swagger Annotation..."
swag fmt

echo "Running swag init..."
swag init
#swag init --exclude pkg/api,pkg/debug
#swag init --exclude pkg/backoffice,pkg/pka,pkg/third_party

if [ $? -ne 0 ]; then
    echo "swag init failed with error code $?."
    exit 1
fi

echo "swag init completed successfully."



echo "Running go run main.go..."
go run main.go 

if [ $? -ne 0 ]; then
    echo "go run main.go failed with error code $?."
    exit 1
fi

