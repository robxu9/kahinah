#!/bin/bash

echo "~> creating build/"
mkdir -p build/

echo "~> generating server binary..."
go build -o build/khserver -i -race -vx

echo "~> generating client files..."
(
  cd client
  ember build -prod -o ../build/public
)

echo "~> done."
