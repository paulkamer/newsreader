#!/bin/bash
cd src/

go test ./... -tags 'excludetest'

go test -tags 'excludetest' -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html
