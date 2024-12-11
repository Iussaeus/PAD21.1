#!/bin/bash

# Start PostgreSQL in the background
docker-entrypoint.sh postgres &

# Run Go application
go run .

